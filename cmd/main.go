package main

import (
	"context"

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
		Port:     8080,
		Env:      *modeFlag,
		Metrics:  true,
		LogLevel: 0,
	}
	app := gtea.NewApp(cfg)
	app.Logger.Info("APP Starting",
		map[string]string{
			"version":   version,
			"gitCommit": buildTime,
			"mode":      *modeFlag,
		})

	h := api.Routes(app)
	ctx := context.Background()
	app.Serve(ctx, h)
}
