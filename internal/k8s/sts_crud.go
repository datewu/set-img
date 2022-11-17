package k8s

import (
	"context"
	"errors"

	"github.com/datewu/gtea/jsonlog"
	apps_v1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newSBio(s *apps_v1.StatefulSet) *Bio {
	bio := &Bio{
		Name: s.Name,
	}
	containes := s.Spec.Template.Spec.Containers
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

// GetSBio get statefulset bio
func GetSBio(ns, name string) (*Bio, error) {
	opts := v1.GetOptions{}
	ctx := context.Background()
	s, err := classicalClientSet.AppsV1().StatefulSets(ns).Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}
	return newSBio(s), nil
}

// ListSts list all statefulset bio in :ns
func ListSts(ns string) ([]*Bio, error) {
	ctx := context.Background()
	opts := v1.ListOptions{}
	stses, err := classicalClientSet.AppsV1().StatefulSets(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	its := stses.Items
	res := make([]*Bio, len(its))
	for i, s := range its {
		res[i] = newSBio(&s)
	}
	return res, nil
}

// SetStsImg update statefulset image field
func SetStsImg(id *ContainerPath) error {
	ctx := context.Background()
	opts := v1.GetOptions{}
	s, err := classicalClientSet.AppsV1().StatefulSets(id.Ns).Get(ctx, id.Name, opts)
	if err != nil {
		jsonlog.Err(err, map[string]string{"sts": id.Name, "msg": "get sts failed"})
		return err
	}
	cpy := s.DeepCopy()
	found := false
	for i, c := range cpy.Spec.Template.Spec.Containers {
		if c.Name == id.CName {
			jsonlog.Info("got new image", map[string]string{"sts": id.Name, "newImg": id.Img})
			cpy.Spec.Template.Spec.Containers[i].Image = id.Img
			found = true
			break
		}
	}
	if !found {
		fErr := errors.New("cannot find container")
		jsonlog.Err(fErr, map[string]string{"sts": id.Name, "container": id.CName, "msg": "cannot find container"})
		return fErr
	}
	uOpts := v1.UpdateOptions{}
	_, err = classicalClientSet.AppsV1().StatefulSets(id.Ns).Update(ctx, cpy, uOpts)
	if err != nil {
		jsonlog.Err(err, map[string]string{"sts": id.Name, "msg": "update sts failed"})
	}
	return err
}
