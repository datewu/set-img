package main

import (
	"flag"
)

var (
	modeFlag   = flag.String("mode", "dev", "runing mode")
	kubeconfig = flag.String("kubeconfig", "in-cluster", "path to kubernetes config file")
)

func parseFlag() {
	flag.Parse()
}
