package main

import (
	"github.com/datewu/gtea"
	"github.com/datewu/set-img/cmd/api"
)

var (
	version   = "1.0.0"
	buildTime string
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
			"version":   version,
			"gitCommit": buildTime,
			"mode":      *modeFlag,
		})

	h := api.Routes(app)
	app.Serve(h)
}
