package main

import (
	"github.com/datewu/gtea"
	"github.com/datewu/set-img/api"
)

func main() {
	parseFlag()
	panicIfErr(initKey)
	panicIfErr(initK8s)
	cfg := &gtea.Config{
		Port:    8080,
		Env:     *modeFlag,
		Metrics: true,
	}
	app := gtea.NewApp(cfg)
	app.Logger.PrintInfo("APP Starting",
		map[string]string{
			"version":   SemVer,
			"gitCommit": GitCommit,
			"mode":      *modeFlag,
		})

	h := api.Routes(app)
	app.Serve(h)
}
