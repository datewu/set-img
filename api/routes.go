package api

import (
	"net/http"

	"github.com/datewu/gtea"
	"github.com/datewu/toushi"
)

func Routes(app *gtea.App) http.Handler {
	r := toushi.New(toushi.DefaultConf())
	addBusinessRoutes(app, r)
	return r.Routes()
}

func addBusinessRoutes(app *gtea.App, r *toushi.Router) {
	th := &tokenHandler{app: app}
	kh := &k8sHandler{app: app}
	r.Get("/api/v1/token", th.getToken)

	r.Get("/api/v1/auth/ping", checkAuth(th.authPing))
	r.Get("/api/v1/auth/list/:ns", checkAuth(kh.listDemo))
	r.Get("/api/v1/auth/get/:ns/:name", checkAuth(kh.getDeployBio))
	r.Post("/api/v1/auth/setdeploy/:ns/image", checkAuth(kh.setDeployImg))
}
