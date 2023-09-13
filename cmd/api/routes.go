package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/handler/sse"
	"github.com/datewu/gtea/handler/static"
	"github.com/datewu/gtea/router"
)

func New(app *gtea.App) http.Handler {
	r := router.DefaultRoutesGroup()
	fs := static.FS{
		NoDir:   true,
		TryFile: []string{},
		Root:    "front/static",
	}
	r.ServeFSWithGzip("/static", fs)
	r.Get("/", index(app))
	r.Get("/debug/reload", sse.NewHandler(reloadSSE{app: app}))
	r.Get("/version", serverVersion(app))
	loginRoutes(app, r)
	myRoutes(app, r)
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

func myRoutes(app *gtea.App, r *router.RoutesGroup) {
	h := &myHandler{app: app}
	my := r.Group("/my", h.auth)
	my.Use(handler.GzipMiddleware)
	my.Delete("/logout", h.logout)
	my.Get("/profile", h.profile)
	my.Get("/deploys", h.deploys)
	my.Get("/sts", h.sts)
	my.Put("/update/resource", h.updateResouce)

}

func addBusinessRoutes(app *gtea.App, r *router.RoutesGroup) {
	kh := &k8sHandler{app: app}
	g := r.Group("/api/v1")
	g.Get("/", showPath)
	a := g.Group("/auth", kh.auth)

	a.Get("/list/:ns/:kind", kh.listBio)
	a.Get("/get/:ns/:kind/:name", kh.getBio)
	a.Post("/setimg", kh.setImg)
}
