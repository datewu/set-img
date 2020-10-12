package k8s

import (
	"path/filepath"

	"github.com/rs/zerolog/log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Conf ...
type Conf struct {
	Mode, ConfFile string
}

var classicalClientSet *kubernetes.Clientset

// InitClientSet ...
func InitClientSet(c *Conf) {
	restConfig := loadRestConfig(c)
	cli, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}
	classicalClientSet = cli
}

func loadRestConfig(conf *Conf) *rest.Config {
	var c *rest.Config
	var err error
	if conf.Mode == "dev" {
		home := homedir.HomeDir()
		conf.ConfFile = filepath.Join(home, ".kube", "config")
	}
	switch conf.ConfFile {
	case "in-cluster":
		log.Debug().Msg("using in-cluster configuration")
		c, err = rest.InClusterConfig()
	default:
		log.Debug().
			Str("config path", conf.ConfFile).
			Msg("using outer cluster configuration")
		c, err = clientcmd.BuildConfigFromFlags("", conf.ConfFile)
	}

	if err != nil {
		log.Panic().Err(err).
			Str("config flag", conf.ConfFile).
			Msg("cannot load kube-client config info")
	}
	return c
}
