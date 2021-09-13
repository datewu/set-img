package k8s

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newDBio(d *apps_v1.Deployment) *Bio {
	bio := &Bio{
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

// GetDBio get deployment Bio
func GetDBio(ns, name string) (*Bio, error) {
	opts := v1.GetOptions{}
	ctx := context.Background()
	d, err := classicalClientSet.AppsV1().Deployments(ns).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	return newDBio(d), nil
}

// ListDeploy list deployment bios
func ListDeploy(ns string) ([]*Bio, error) {
	ctx := context.Background()
	opts := v1.ListOptions{}
	deploys, err := classicalClientSet.AppsV1().Deployments(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	its := deploys.Items
	res := make([]*Bio, len(its))
	for i, d := range its {
		res[i] = newDBio(&d)
	}
	return res, nil
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
