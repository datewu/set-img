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

func (t tokenHandler) getToken(w http.ResponseWriter, r *http.Request) {
	h := toushi.NewHandleHelper(w, r)
	// TODO needs more security
	if !strings.HasPrefix(r.Host, "localhost:") {
		h.BadRequestMsg("route only available to localhost")
		return
	}
	token, err := auth.NewToken()
	if err != nil {
		h.ServerErr(err)
		return
	}
	toushi.OKJSON(w, toushi.Envelope{"token": token})
}

func (t tokenHandler) authPing(w http.ResponseWriter, r *http.Request) {
	msg := "ping from auth, you've been  authenticated"
	toushi.WriteStr(w, http.StatusOK, msg, nil)
}

type k8sHandler struct {
	app *gtea.App
}

func (k k8sHandler) listBio(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	kind := toushi.ReadParams(r, "kind")
	data := toushi.Envelope{}
	h := toushi.NewHandleHelper(w, r)
	switch kind {
	case "deploy":
		ls, err := k8s.ListDeploy(ns)
		if err != nil {
			h.ServerErr(err)
			return
		}
		data["developments"] = ls
	case "sts":
		ls, err := k8s.ListSts(ns)
		if err != nil {
			h.ServerErr(err)
			return
		}
		data["sts"] = ls
	default:
		h.BadRequestMsg("only support deploy/sts two kind resource")
		return
	}
	toushi.OKJSON(w, data)
}

func (k k8sHandler) getBio(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	kind := toushi.ReadParams(r, "kind")
	name := toushi.ReadParams(r, "name")
	data := toushi.Envelope{}
	h := toushi.NewHandleHelper(w, r)
	switch kind {
	case "deploy":
		b, err := k8s.GetDBio(ns, name)
		if err != nil {
			h.ServerErr(err)
			return
		}
		data["bio"] = b
	case "sts":
		b, err := k8s.GetSBio(ns, name)
		if err != nil {
			h.ServerErr(err)
			return
		}
		data["bio"] = b
	default:
		h.BadRequestMsg("only support deploy/sts two kind resource")
		return
	}
	toushi.OKJSON(w, data)
}

func (k k8sHandler) setImg(w http.ResponseWriter, r *http.Request) {
	id := new(k8s.ContainerPath)
	h := toushi.NewHandleHelper(w, r)
	err := toushi.ReadJSON(w, r, id)
	if err != nil {
		h.BadRequestErr(err)
		return
	}
	switch id.Kind {
	case "deploy":
		err = k8s.SetDeployImg(id)
	case "sts":
		err = k8s.SetStsImg(id)
	default:
		h.BadRequestMsg("only support deploy/sts two kind resource")
		return
	}
	if err != nil {
		h.ServerErr(err)
		return
	}
	toushi.OKJSON(w, toushi.Envelope{"payload": id})
}
