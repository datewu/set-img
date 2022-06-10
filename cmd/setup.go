package main

import (
	"github.com/datewu/set-img/internal/k8s"
)

func initK8s() error {
	k8sConf := &k8s.Conf{
		Mode:     env,
		ConfFile: kubeconfig,
	}
	return k8s.InitClientSet(k8sConf)
}
