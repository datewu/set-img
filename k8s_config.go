package main

import (
	"path/filepath"

	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var classicalClientSet = newClientSet()

func newClientSet() *kubernetes.Clientset {
	restConfig := loadRestConfig()
	cli, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	return cli
}

func loadRestConfig() *rest.Config {
	var c *rest.Config
	var err error
	k := *kubeconfig
	if *modeFlag == "dev" {
		home := homedir.HomeDir()
		k = filepath.Join(home, ".kube", "config")
	}
	switch k {
	case "in-cluster":
		log.Debug().Msg("using in-cluster configuration")
		c, err = rest.InClusterConfig()
	default:
		log.Debug().
			Str("config path", k).
			Msg("using outer cluster configuration")
		c, err = clientcmd.BuildConfigFromFlags("", k)
	}

	if err != nil {
		log.Panic().Err(err).
			Str("config flag", k).
			Msg("cannot load kube-client config info")
	}
	return c
}
