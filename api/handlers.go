package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/datewu/gtea"
	"github.com/datewu/set-img/auth"
	"github.com/datewu/set-img/k8s"
	"github.com/datewu/toushi"
)

type tokenHandler struct {
	app *gtea.App
}

func (h tokenHandler) getToken(w http.ResponseWriter, r *http.Request) {
	// TODO needs more security
	if !strings.HasPrefix(r.Host, "localhost:") {
		err := errors.New("route only available to localhost")
		toushi.BadRequestResponse(err)(w, r)
		return
	}
	token, err := auth.NewToken()
	if err != nil {
		toushi.ServerErrResponse(err)(w, r)
		return
	}
	toushi.OKJSON(w, r, toushi.Envelope{"token": token})
}

func (h tokenHandler) authPing(w http.ResponseWriter, r *http.Request) {
	msg := "ping from auth, you've been  authenticated"
	toushi.WriteStr(w, http.StatusOK, msg, nil)
}

type k8sHandler struct {
	app *gtea.App
}

func (h k8sHandler) listDemo(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	ls, err := k8s.ListDemo(ns)
	if err != nil {
		toushi.ServerErrResponse(err)(w, r)
		return
	}
	data := toushi.Envelope{"developments": ls}
	toushi.OKJSON(w, r, data)
}

func (h k8sHandler) getDeployBio(w http.ResponseWriter, r *http.Request) {
	ns := toushi.ReadParams(r, "ns")
	name := toushi.ReadParams(r, "name")
	b, err := k8s.GetDBio(ns, name)
	if err != nil {
		toushi.ServerErrResponse(err)(w, r)
		return
	}
	data := toushi.Envelope{"bio": b}
	toushi.OKJSON(w, r, data)
}
func (h k8sHandler) setDeployImg(w http.ResponseWriter, r *http.Request) {

	id := new(k8s.ContainerPath)
	err := toushi.ReadJSON(w, r, id)
	if err != nil {
		toushi.BadRequestResponse(err)(w, r)
		return
	}
	id.Ns = toushi.ReadParams(r, "ns")
	err = k8s.SetDeployImg(id)
	if err != nil {
		toushi.ServerErrResponse(err)(w, r)
		return
	}
	toushi.OKJSON(w, r, toushi.Envelope{"payload": id})
}
