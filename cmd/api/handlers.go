package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/set-img/front"
	"github.com/datewu/set-img/internal/auth/github"
	"github.com/datewu/set-img/internal/k8s"
)

func index(a *gtea.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		view := front.IndexView{}
		layout := front.NewLayout("", a.Env())
		if layout.Env == gtea.DevEnv {
			layout.User = "datewu"
			if err := view.FullPageRender(w, layout); err != nil {
				handler.ServerErr(w, err)
			}
			return
		}
		token, err := r.Cookie(github.CookieName)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				jsonlog.Info("no cookie found")
			} else {
				jsonlog.Err(err)
			}
			if err := view.FullPageRender(w, layout); err != nil {
				handler.ServerErr(w, err)
			}
			return
		}
		user, err := github.GetUser(token.Value)
		if err != nil {
			handler.ClearSimpleCookie(w, github.CookieName)
			if err := view.FullPageRender(w, layout); err != nil {
				handler.ServerErr(w, err)
			}
			return
		}
		layout.User = user.Login
		if err := view.FullPageRender(w, layout); err != nil {
			handler.ServerErr(w, err)
		}
	}
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

type k8sHandler struct {
	app  *gtea.App
	user string
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

	go func() {
		site := handler.ReadQuery(r, "site", "")
		if len(site) < 3 {
			fmt.Printf("ignore site: %q \n", site)
			return
		}
		time.Sleep(30 * time.Second)
		set_cdn(site)

	}()
	handler.OKJSON(w, handler.Envelope{"payload": id})
}

// set poor man's cdn
// site parameter should contain domian only, no schema included
func set_cdn(site string) {
	type Site struct {
		Origin string `json:"origin"`
	}

	url := "http://r2-s3/auto-origin"

	fmt.Printf("start handle site: %q \n", url)
	defer func() {
		fmt.Printf("handle site: %q done\n", url)
	}()

	data := Site{
		Origin: site,
	}

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		// fmt.Errorf("failed to marshal JSON payload: %w \n", err)
		fmt.Println("failed to marshal JSON payload ")
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		// fmt.Errorf("failed to send HTTP POST request: %w \n", err)
		fmt.Println("failed to send HTTP POST request")
		return
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Errorf("failed to read response body: %w \n", err)
		fmt.Println("failed to read response body")
		return
	}
	fmt.Printf("set poor man's cdn 'ok', %s \n", body)

}
