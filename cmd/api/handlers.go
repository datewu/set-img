package api

import (
	"fmt"
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/set-img/internal/k8s"
)

func showPath(w http.ResponseWriter, r *http.Request) {
	url := `https://%s/api/v1/auth/setimg`
	curl := `curl -XPOST $URL -H 'Authorization: $TOKEN' \
	--data-binary '{"namespace":"CHANGE-ME","kind": "CHANGE-ME-deploy/sts","name":"CHANGE-ME","container_name":"img","img":"${{ steps.prep.outputs.tags }}"}'`
	uri := fmt.Sprintf(url, r.Host)
	data := handler.Envelope{
		"url":  uri,
		"curl": curl,
	}
	handler.OKJSON(w, data)
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
		ls, err := k8s.ListDeploy(ns)
		if err != nil {
			handler.ServerErr(w, err)
			return
		}
		data["developments"] = ls
	case "sts":
		ls, err := k8s.ListSts(ns)
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

func (k k8sHandler) getBio(w http.ResponseWriter, r *http.Request) {
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
