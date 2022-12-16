package main

import (
	"context"
	"flag"

	"github.com/datewu/gtea"
	"github.com/datewu/gtea/jsonlog"
	"github.com/datewu/gtea/utils"
	"github.com/datewu/set-img/cmd/api"
)

var (
	version   = "1.0.0"
	buildTime string
)
var (
	port       int
	env        string
	kubeconfig string
)

func main() {

	flag.IntVar(&port, "port", 8080, "API server port")
	flag.StringVar(&env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&kubeconfig, "kubeconfig", "in-cluster", "path to kubernetes config file")

	flag.Parse()

	utils.PanicFn(initK8s)

	cfg := &gtea.Config{
		Port:     port,
		Env:      env,
		Metrics:  true,
		LogLevel: jsonlog.LevelInfo,
	}
	ctx := context.Background()
	app := gtea.NewApp(ctx, cfg)
	app.Logger.Info("APP Starting",
		map[string]any{
			"version":   version,
			"gitCommit": buildTime,
			"mode":      env,
		})
	app.AddMetaData("version", version)
	app.Serve(ctx, api.New(app))
}
