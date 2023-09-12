package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/set-img/front"
	"github.com/datewu/set-img/internal/k8s"
)

func index(a *gtea.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		view := front.IndexView{}
		if a.Env() == gtea.DevEnv {
			view.User = "datewu"
			view.Render(w)
			return
		}
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
