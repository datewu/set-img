package k8s

import (
	"context"

	"github.com/rs/zerolog/log"
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

// ContainerPath ...
type ContainerPath struct {
	Ns    string
	Name  string `json:"deploy_name" binding:"required"`
	CName string `json:"container_name" binding:"required"`
	Img   string `json:"img" binding:"required"`
}

// SetDeployImg ...
func SetDeployImg(id *ContainerPath) {
	ctx := context.Background()
	opts := v1.GetOptions{}
	d, err := classicalClientSet.AppsV1().Deployments(id.Ns).Get(ctx, id.Name, opts)
	if err != nil {
		log.Error().Err(err).
			Str("name", id.Name).
			Msg("get deploy failed")
		return
	}
	cpy := d.DeepCopy()
	found := false
	for _, c := range cpy.Spec.Template.Spec.Containers {
		if c.Name == id.CName {
			c.Image = id.Img
			found = true
			break
		}
	}
	if !found {
		log.Error().Err(err).
			Str("deploy", id.Name).
			Str("container", id.CName).
			Msg("canot find container")
		return
	}
	uOpts := v1.UpdateOptions{}
	_, err = classicalClientSet.AppsV1().Deployments(id.Ns).Update(ctx, cpy, uOpts)
	if err != nil {
		log.Error().Err(err).
			Msg("update deploy failed")
	}
}
