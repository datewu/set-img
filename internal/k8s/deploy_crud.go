package k8s

import (
	"context"
	"errors"

	"github.com/datewu/gtea/jsonlog"
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
		jsonlog.Err(err, map[string]string{"name": id.Name, "msg": "get deploy failed"})
		return err
	}
	cpy := d.DeepCopy()
	found := false
	for i, c := range cpy.Spec.Template.Spec.Containers {
		if c.Name == id.CName {
			jsonlog.Info("got new image tag", map[string]string{"deploy": id.Name, "image": id.Img})
			cpy.Spec.Template.Spec.Containers[i].Image = id.Img
			found = true
			break
		}
	}
	if !found {
		fErr := errors.New("cannot find container")
		jsonlog.Err(fErr, map[string]string{"deploy": id.Name, "image": id.Img, "container": id.CName})
		return fErr
	}
	uOpts := v1.UpdateOptions{}
	_, err = classicalClientSet.AppsV1().Deployments(id.Ns).Update(ctx, cpy, uOpts)
	if err != nil {
		jsonlog.Err(err, map[string]string{"deploy": id.Name, "image": id.Img, "msg": "update deploy failed"})
	}
	return err
}
