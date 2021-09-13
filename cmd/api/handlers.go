package api

import (
	"net/http"
	"strings"

	"github.com/datewu/gtea"
	"github.com/datewu/set-img/internal/auth"
	"github.com/datewu/set-img/internal/k8s"
	"github.com/datewu/toushi"
)

type tokenHandler struct {
	app *gtea.App
}

func (h tokenHandler) getToken(w http.ResponseWriter, r *http.Request) {
	// TODO needs more security
	if !strings.HasPrefix(r.Host, "localhost:") {
		toushi.HandleBadRequestMsg("route only available to localhost")(w, r)
		return
	}
	token, err := auth.NewToken()
	if err != nil {
		toushi.HandleServerErr(err)(w, r)
		return
	}
	toushi.OKJSON(w, toushi.Envelope{"token": token})
}

func (h tokenHandler) authPing(w http.ResponseWriter, r *http.Request) {
	msg := "ping from auth, you've been  authenticated"
	toushi.WriteStr(w, http.StatusOK, msg, nil)
}

type k8sHandler struct {
	app *gtea.App
}

func (h k8sHandler) listBio(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	kind := toushi.ReadParams(r, "kind")
	data := toushi.Envelope{}
	switch kind {
	case "deploy":
		ls, err := k8s.ListDeploy(ns)
		if err != nil {
			toushi.HandleServerErr(err)(w, r)
			return
		}
		data["developments"] = ls
	case "sts":
		ls, err := k8s.ListSts(ns)
		if err != nil {
			toushi.HandleServerErr(err)(w, r)
			return
		}
		data["sts"] = ls
	default:
		toushi.HandleBadRequestMsg("only support deploy/sts two kind resource")(w, r)
		return
	}
	toushi.OKJSON(w, data)
}

func (h k8sHandler) getBio(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	kind := toushi.ReadParams(r, "kind")
	name := toushi.ReadParams(r, "name")
	data := toushi.Envelope{}
	switch kind {
	case "deploy":
		b, err := k8s.GetDBio(ns, name)
		if err != nil {
			toushi.HandleServerErr(err)(w, r)
			return
		}
		data["bio"] = b
	case "sts":
		b, err := k8s.GetSBio(ns, name)
		if err != nil {
			toushi.HandleServerErr(err)(w, r)
			return
		}
		data["bio"] = b
	default:
		toushi.HandleBadRequestMsg("only support deploy/sts two kind resource")(w, r)
		return
	}
	toushi.OKJSON(w, data)
}

func (h k8sHandler) setImg(w http.ResponseWriter, r *http.Request) {
	id := new(k8s.ContainerPath)
	err := toushi.ReadJSON(w, r, id)
	if err != nil {
		toushi.HandleBadRequestErr(err)(w, r)
		return
	}
	switch id.Kind {
	case "deploy":
		err = k8s.SetDeployImg(id)
	case "sts":
		err = k8s.SetStsImg(id)
	default:
		toushi.HandleBadRequestMsg("only support deploy/sts two kind resource")(w, r)
		return
	}
	if err != nil {
		toushi.HandleServerErr(err)(w, r)
		return
	}
	toushi.OKJSON(w, toushi.Envelope{"payload": id})
}
