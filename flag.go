package main

import (
	"flag"

	"github.com/rs/zerolog/log"
)

var (
	modeFlag   = flag.String("mode", "dev", "runing mode")
	kubeconfig = flag.String("kubeconfig", "in-cluster", "path to kubernetes config file")
)

func parseFlag() {
	flag.Parse()
	log.Info().
		Str("version", SemVer).
		Str("gitCommit", GitCommit).
		Msg("APP starting ...")

	log.Info().
		Str("mode", *modeFlag).
		Msg("APP arguments")
}
