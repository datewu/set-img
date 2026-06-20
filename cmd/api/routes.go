package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/handler"
	"github.com/datewu/gtea/handler/sse"
	"github.com/datewu/gtea/handler/static"
	"github.com/datewu/gtea/router"
	"github.com/datewu/set-img/front"
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
	if app.Env() == gtea.DevEnv {
		r.Get("/debug/reload",
			sse.NewHandler(newReloadSSE(app, front.InitOrReload,
				"front", "front/static")))
	}
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
	my.Get("/logs", h.logs)

	myGzip := my.Group("", handler.GzipMiddleware)
	myGzip.Delete("/logout", h.logout)
	myGzip.Get("/profile", h.profile)
	myGzip.Get("/deploys", h.deploys)
	myGzip.Get("/sts", h.sts)
	myGzip.Put("/update/resource", h.updateResouce)
	myGzip.Post("/update/cdn", h.updateCDN)
	myGzip.Get("/pods", h.listPods)
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