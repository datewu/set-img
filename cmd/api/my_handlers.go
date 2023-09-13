package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/set-img/front"
	"github.com/datewu/set-img/internal/auth/github"
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

type ghLoginHandler struct {
	app *gtea.App
}

func (ghLoginHandler) init(w http.ResponseWriter, r *http.Request) {
	htmx := fmt.Sprintf(`<a hx-boost="false" href="%s/login/oauth/authorize?client_id=%s">Login</a>`,
		github.URL, os.Getenv("GITHUB-APP-CID"))
	handler.OKText(w, htmx)
}

func (g ghLoginHandler) callback(w http.ResponseWriter, r *http.Request) {
	code := handler.ReadQuery(r, "code", "")
	if code == "" {
		handler.BadRequestErr(w, fmt.Errorf("code is empty"))
		return
	}
	token, err := github.GetToken(os.Getenv("GITHUB-APP-CID"),
		os.Getenv("GITHUB-APP-SECRET"), code)
	if err != nil {
		// handler.BadRequestErr(w, err)
		handler.ServerErr(w, err)
		return
	}
	handler.SetSimpleCookie(w, r, github.CookieName, token.AccessToken)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

type myHandler struct {
	app   *gtea.App
	user  string
	token string
}

func (m *myHandler) profile(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "" {
		htmx := fmt.Sprintf(`<span> hello %s</span>`, m.user)
		handler.OKText(w, htmx)
		return
	}
	view := front.ProfileView{User: m.user}
	if err := view.Render(w); err != nil {
		handler.ServerErr(w, err)
	}
}

func (m *myHandler) logout(w http.ResponseWriter, r *http.Request) {
	handler.ClearSimpleCookie(w, github.CookieName)
	// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	w.Header().Add("HX-Redirect", "/")
	handler.OKText(w, "logout")
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
	if r.Header.Get("HX-Request") == "" {
		if err := view.Render(w, m.user); err != nil {
			handler.ServerErr(w, err)
		}
		return
	}
	if err := view.Render(w, ""); err != nil {
		handler.ServerErr(w, err)
	}
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
	if r.Header.Get("HX-Request") == "" {
		if err := view.Render(w, m.user); err != nil {
			handler.ServerErr(w, err)
		}
		return
	}
	if err := view.Render(w, ""); err != nil {
		handler.ServerErr(w, err)
	}
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
