package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
)

func New(app *gtea.App) http.Handler {
	r := handler.DefaultRouterGroup()
	addBusinessRoutes(app, r)
	return r
}

func addBusinessRoutes(app *gtea.App, r *handler.RouterGroup) {
	th := &tokenHandler{app: app}
	kh := &k8sHandler{app: app}
	g := r.Group("/api/v1")
	g.Get("/token", th.getToken)
	a := g.Group("/auth", checkAuth)

	a.Get("/ping", th.authPing)
	a.Get("/list/:ns/:kind", kh.listBio)
	a.Get("/get/:ns/:kind/:name", kh.getBio)
	a.Post("/setimg", kh.setImg)
}
