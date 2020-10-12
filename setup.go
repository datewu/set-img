package main

import (
	"github.com/datewu/set-img/api"
	"github.com/datewu/set-img/k8s"
)

func initK8s() {
	k8sConf := &k8s.Conf{
		Mode:     *modeFlag,
		ConfFile: *kubeconfig,
	}
	k8s.InitClientSet(k8sConf)
}

func server() {
	apiConf := &api.Conf{
		Mode: *modeFlag,
		Addr: ":8080",
	}

	api.Server(apiConf)
}
