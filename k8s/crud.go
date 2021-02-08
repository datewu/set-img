package k8s

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DBio ...
type DBio struct {
	Name       string    `json:"name"`
	Containers []*ConBio `json:"containers"`
}

func newDBio(d *apps_v1.Deployment) *DBio {
	bio := &DBio{
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
	bio.Containers = cs
	return bio
}

// ConBio ..
type ConBio struct {
	Name  string `json:"name"`
	Image string `json:"img"`
	Pull  string `json:"pull"`
}

// GetDBio by name
func GetDBio(ns, name string) (*DBio, error) {
	opts := v1.GetOptions{}
	ctx := context.Background()
	d, err := classicalClientSet.AppsV1().Deployments(ns).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	return newDBio(d), nil
}

// ListDemo ...
func ListDemo(ns string) ([]*DBio, error) {
	ctx := context.Background()
	opts := v1.ListOptions{}
	deploys, err := classicalClientSet.AppsV1().Deployments(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	its := deploys.Items
	res := make([]*DBio, len(its))
	for i, d := range its {
		res[i] = newDBio(&d)
	}
	return res, nil
}

// ContainerPath ...
type ContainerPath struct {
	Ns    string
	Name  string `json:"deploy_name" binding:"required"`
	CName string `json:"container_name" binding:"required"`
	Img   string `json:"img" binding:"required"`
}

func (c *ContainerPath) formatImg() error {
	const prefix = "refs/tags/"
	a := strings.Split(c.Img, ":")
	if len(a) != 2 {
		return errors.New("invalid imgage format")
	}
	if !strings.HasPrefix(a[1], prefix) {
		return nil
		//return errors.New("invalid github tag")
	}
	c.Img = a[0] + ":" + a[1][len(prefix):]
	return nil
}

// SetDeployImg ...
func SetDeployImg(id *ContainerPath) error {
	ctx := context.Background()
	opts := v1.GetOptions{}
	d, err := classicalClientSet.AppsV1().Deployments(id.Ns).Get(ctx, id.Name, opts)
	if err != nil {
		log.Error().Err(err).
			Str("name", id.Name).
			Msg("get deploy failed")
		return err
	}
	cpy := d.DeepCopy()
	found := false
	for i, c := range cpy.Spec.Template.Spec.Containers {
		if c.Name == id.CName {
			err := id.formatImg()
			if err != nil {
				return err
			}
			log.Info().
				Str("deploy", id.Name).
				Str("newImg", id.Img).
				Msg("got new image")
			cpy.Spec.Template.Spec.Containers[i].Image = id.Img
			found = true
			break
		}
	}
	if !found {
		log.Error().Err(err).
			Str("deploy", id.Name).
			Str("container", id.CName).
			Msg("canot find container")
		return errors.New("cannot find container")
	}
	uOpts := v1.UpdateOptions{}
	_, err = classicalClientSet.AppsV1().Deployments(id.Ns).Update(ctx, cpy, uOpts)
	if err != nil {
		log.Error().Err(err).
			Msg("update deploy failed")
	}
	return err
}
