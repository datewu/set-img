package main

import (
	"github.com/datewu/set-img/internal/k8s"
)

func panicIfErr(fn func() error) {
	err := fn()
	if err != nil {
		panic(err)
	}
}

func initK8s() error {
	k8sConf := &k8s.Conf{
		Mode:     env,
		ConfFile: kubeconfig,
	}
	return k8s.InitClientSet(k8sConf)
}
