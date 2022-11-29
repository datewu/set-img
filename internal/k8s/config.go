package k8s

import (
	"path/filepath"

	"github.com/datewu/gtea/jsonlog"
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
func InitClientSet(c *Conf) error {
	restConfig := loadRestConfig(c)
	cli, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	classicalClientSet = cli
	return nil
}

func loadRestConfig(conf *Conf) *rest.Config {
	var c *rest.Config
	var err error
	if conf.Mode == "development" {
		home := homedir.HomeDir()
		conf.ConfFile = filepath.Join(home, ".kube", "config")
	}
	switch conf.ConfFile {
	case "in-cluster":
		jsonlog.Debug("using in-cluster configuration", nil)
		c, err = rest.InClusterConfig()
	default:
		jsonlog.Debug("using outer cluster configuration", map[string]interface{}{"configPath": conf.ConfFile})
		c, err = clientcmd.BuildConfigFromFlags("", conf.ConfFile)
	}

	if err != nil {
		jsonlog.Err(err, map[string]interface{}{"configPath": conf.ConfFile, "msg": "cannot load kube-client config info"})
		panic(err)
	}
	return c
}
