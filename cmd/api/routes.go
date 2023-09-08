package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/handler/static"
	"github.com/datewu/gtea/router"
)

func New(app *gtea.App) http.Handler {
	r := router.DefaultRoutesGroup()
	fs := static.FS{
		NoDir:   true,
		TryFile: []string{},
		Root:    "front",
	}
	r.ServeFSWithGzip("/", fs)
	r.Get("/version", serverVersion(app))
	loginRoutes(app, r)
	addBusinessRoutes(app, r)
	return r.Handler()
}

func loginRoutes(app *gtea.App, r *router.RoutesGroup) {
	login := r.Group("/login")
	gh := login.Group("/github")
	g := &ghLoginHandler{app: app}
	gh.Get("/init", g.init)
	gh.Get("/callback", g.callback)

}

func addBusinessRoutes(app *gtea.App, r *router.RoutesGroup) {
	th := &tokenHandler{app: app}
	kh := &k8sHandler{app: app}
	g := r.Group("/api/v1")
	g.Get("/", showPath)
	a := g.Group("/auth", handler.TokenMiddleware(checkAuth))

	a.Get("/ping", th.authPing)
	a.Get("/list/:ns/:kind", kh.listBio)
	a.Get("/get/:ns/:kind/:name", kh.getBio)
	a.Post("/setimg", kh.setImg)
}
