package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/set-img/front"
	"github.com/datewu/set-img/internal/k8s"
)

func serverVersion(a *gtea.App) func(w http.ResponseWriter, r *http.Request) {
	version := a.GetMetaData("version")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") == "application/json" || r.Header.Get("HX-Request") == "" {
			handler.Version(version)
			return
		}
		htmx := fmt.Sprintf(`<span>%s</span>`, version)
		handler.OKText(w, htmx)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	view := front.IndexView{}
	token, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			jsonlog.Info("no cookie found")
		}
		jsonlog.Err(err)
		view.Render(w)
		return
	}
	//func (ghLoginHandler) userInfo(token string) (*UserInfo, error) {
	g := ghLoginHandler{}
	user, err := g.userInfo(token.Value)
	if err != nil {
		view.Render(w)
		return
	}
	view.User = user.Login
	view.Render(w)
}

type curlCmd struct {
	Method     string            `json:"method"`
	Header     map[string]string `json:"headers"`
	URL        string            `json:"url"`
	BinaryData map[string]any    `json:"data"`
}

func showPath(w http.ResponseWriter, r *http.Request) {
	url := `https://%s/api/v1/auth/setimg`
	uri := fmt.Sprintf(url, r.Host)
	curl := curlCmd{
		URL:    uri,
		Method: "POST",
		Header: map[string]string{"Authorization": "$TOKEN"},
		BinaryData: map[string]any{
			"namespace":      "CHANGE-ME",
			"kind":           "CHANGE-ME-TO-deploy/sts",
			"name":           "CHANGE-ME",
			"container_name": "CHANGE-ME",
			"img":            "${{ steps.prep.outputs.tags }}",
		},
	}
	handler.OKJSON(w, curl)
}

type tokenHandler struct {
	app *gtea.App
}

func (t tokenHandler) authPing(w http.ResponseWriter, r *http.Request) {
	msg := "ping from auth, you've been  authenticated"
	handler.WriteStr(w, http.StatusOK, msg, nil)
}

type k8sHandler struct {
	app *gtea.App
}

func (k k8sHandler) listBio(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadPathParam(r, "ns")
	kind := handler.ReadPathParam(r, "kind")
	data := handler.Envelope{}
	switch kind {
	case "deploy":
		ls, err := k8s.ListBios(ns)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		data["developments"] = ls
	case "sts":
		ls, err := k8s.ListStsBios(ns)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		data["sts"] = ls
	default:
		handler.BadRequestMsg(w, "only support deploy/sts two kind resource")
		return
	}
	handler.OKJSON(w, data)
}

func (k8sHandler) getBio(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadPathParam(r, "ns")
	kind := handler.ReadPathParam(r, "kind")
	name := handler.ReadPathParam(r, "name")
	data := handler.Envelope{}
	switch kind {
	case "deploy":
		b, err := k8s.GetDBio(ns, name)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		data["bio"] = b
	case "sts":
		b, err := k8s.GetSBio(ns, name)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		data["bio"] = b
	default:
		handler.BadRequestMsg(w, "only support deploy/sts two kind resource")
		return
	}
	handler.OKJSON(w, data)
}

func (k k8sHandler) setImg(w http.ResponseWriter, r *http.Request) {
	id := new(k8s.ContainerPath)
	err := handler.ReadJSON(r, id)
	if err != nil {
		handler.BadRequestErr(w, err)
		return
	}
	switch id.Kind {
	case "deploy":
		err = k8s.SetDeployImg(id)
	case "sts":
		err = k8s.SetStsImg(id)
	default:
		handler.BadRequestMsg(w, "only support deploy/sts two kind resource")
		return
	}
	if err != nil {
		handler.ServerErr(w, err)
		return
	}
	handler.OKJSON(w, handler.Envelope{"payload": id})
}

type ghLoginHandler struct {
	app *gtea.App
}

// UserInfo is the github user info
type UserInfo struct {
	Login             string    `json:"login"`
	ID                int       `json:"id"`
	NodeID            string    `json:"node_id"`
	AvatarURL         string    `json:"avatar_url"`
	GravatarID        string    `json:"gravatar_id"`
	URL               string    `json:"url"`
	HTMLURL           string    `json:"html_url"`
	FollowersURL      string    `json:"followers_url"`
	FollowingURL      string    `json:"following_url"`
	GistsURL          string    `json:"gists_url"`
	StarredURL        string    `json:"starred_url"`
	SubscriptionsURL  string    `json:"subscriptions_url"`
	OrganizationsURL  string    `json:"organizations_url"`
	ReposURL          string    `json:"repos_url"`
	EventsURL         string    `json:"events_url"`
	ReceivedEventsURL string    `json:"received_events_url"`
	Type              string    `json:"type"`
	SiteAdmin         bool      `json:"site_admin"`
	Name              string    `json:"name"`
	Company           string    `json:"company"`
	Blog              string    `json:"blog"`
	Location          string    `json:"location"`
	Email             string    `json:"email"`
	Hireable          bool      `json:"hireable"`
	Bio               string    `json:"bio"`
	TwitterUsername   string    `json:"twitter_username"`
	PublicRepos       int       `json:"public_repos"`
	PublicGists       int       `json:"public_gists"`
	Followers         int       `json:"followers"`
	Following         int       `json:"following"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// GhToken is the github token
type GhToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (ghLoginHandler) getToken(code string) (*GhToken, error) {
	data := fmt.Sprintf(`{"client_id":"%s","client_secret":"%s","code":"%s"}`,
		os.Getenv("GITHUB-APP-CID"), os.Getenv("GITHUB-APP-SECRET"), code)
	reader := strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-GitHub-OTP", "")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var token GhToken
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (ghLoginHandler) userInfo(token string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var info UserInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (ghLoginHandler) init(w http.ResponseWriter, r *http.Request) {
	htmx := fmt.Sprintf(`
	<a hx-boost="false" href="https://github.com/login/oauth/authorize?client_id=%s">Login with GitHub</a>
	`, os.Getenv("GITHUB-APP-CID"))
	handler.OKText(w, htmx)
}

func (g ghLoginHandler) callback(w http.ResponseWriter, r *http.Request) {
	code := handler.ReadQuery(r, "code", "")
	if code == "" {
		handler.BadRequestErr(w, fmt.Errorf("code is empty"))
		return
	}
	token, err := g.getToken(code)
	if err != nil {
		// handler.BadRequestErr(w, err)
		handler.ServerErr(w, err)
		return
	}
	handler.SetSimpleCookie(w, r, "access_token", token.AccessToken)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type myHandler struct {
	app   *gtea.App
	user  string
	token string
}

func (m *myHandler) middlerware(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		co, err := r.Cookie("access_token")
		if err != nil {
			handler.BadRequestMsg(w, "missing access_token cookie")
			return
		}
		t := co.Value
		gh := ghLoginHandler{}
		user, err := gh.userInfo(t)
		if err != nil {
			handler.ClearSimpleCookie(w, "access_token")
			handler.ServerErr(w, err)
			return
		}
		m.user = user.Login
		m.token = t
		next(w, r)
	}
	return middle
}

func (m *myHandler) profile(w http.ResponseWriter, r *http.Request) {
	htmx := fmt.Sprintf(`<span> hello %s</span>`, m.user)
	handler.OKText(w, htmx)
}

func (m *myHandler) deploys(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadQuery(r, "ns", "wu")
	view := &front.TableView{
		Description: "deployments by " + m.user,
		Namespace:   ns,
		Kind:        "deploy",
	}
	label := fmt.Sprintf("image-user=%s", m.user)
	ds, err := k8s.ListDeployWithLabels(ns, label)
	if err != nil {
		handler.ServerErr(w, err)
		return
	}
	view.AddDeploys(ds)
	view.Render(w)
}

func (m *myHandler) sts(w http.ResponseWriter, r *http.Request) {
	ns := handler.ReadQuery(r, "ns", "wu")
	view := &front.TableView{
		Description: "statefulsets by " + m.user,
		Namespace:   ns,
		Kind:        "sts",
	}
	label := fmt.Sprintf("image-user=%s", m.user)
	ss, err := k8s.ListStsWithLabels(ns, label)
	if err != nil {
		handler.ServerErr(w, err)
		return
	}
	view.AddSts(ss)
	view.Render(w)
}

func (m *myHandler) updateResouce(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handler.BadRequestErr(w, err)
		return
	}
	c := k8s.ContainerPath{
		Ns:    r.FormValue("ns"),
		Kind:  r.FormValue("kind"),
		Name:  r.FormValue("name"),
		CName: r.FormValue("cname"),
		Img:   r.FormValue("image"),
	}
	err = c.UpdateResource(fmt.Sprintf("image-user=%s", m.user))
	if err != nil {
		handler.ServerErr(w, err)
		return
	}
	handler.OKText(w, "ok")
}
