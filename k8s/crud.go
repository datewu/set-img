package k8s

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DBio ...
type DBio struct {
	Name       string    `json:"name"`
	Containers []*ConBio `json:"containers"`
}

// ConBio ..
type ConBio struct {
	Name  string `json:"name"`
	Image string `json:"img"`
	Pull  string `json:"pull"`
}

// ListDemo ...
func ListDemo(ns string) []*DBio {
	ctx := context.Background()
	opts := v1.ListOptions{}
	deploys, err := classicalClientSet.AppsV1().Deployments(ns).List(ctx, opts)
	if err != nil {
		panic(err)
	}
	its := deploys.Items
	res := make([]*DBio, len(its))
	for i, d := range its {
		de := &DBio{
			Name: d.Name,
		}
		containes := d.Spec.Template.Spec.Containers
		cs := make([]*ConBio, len(containes))
		for i, c := range containes {
			cs[i] = &ConBio{
				Name:  c.Name,
				Image: c.Image,
				Pull:  string(c.ImagePullPolicy),
			}
		}
		de.Containers = cs
		res[i] = de
	}
	return res
}

func updateDeploy() {

}
