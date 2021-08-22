package main

import (
	"path"

	"github.com/datewu/set-img/auth"
	"github.com/datewu/set-img/k8s"
)

func panicIfErr(fn func() error) {
	err := fn()
	if err != nil {
		panic(err)
	}
}

func initKey() error {
	fn := "private_key_for_sign.pem"
	if *modeFlag == "production" {
		fn = path.Join("/opt", fn)
	}
	return auth.InitKeys(fn)
}

func initK8s() error {
	k8sConf := &k8s.Conf{
		Mode:     *modeFlag,
		ConfFile: *kubeconfig,
	}
	return k8s.InitClientSet(k8sConf)
}
