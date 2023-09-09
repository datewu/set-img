package k8s

import (
	"context"
	"errors"
	"time"

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

// ListStsWithLables list statefulset
func ListStsWithLabels(ns, label string) ([]apps_v1.StatefulSet, error) {
	ctx := context.Background()
	opts := v1.ListOptions{LabelSelector: label}
	stses, err := classicalClientSet.AppsV1().StatefulSets(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return stses.Items, nil
}

// ListStsBios list all statefulset bio in :ns
func ListStsBios(ns string) ([]*Bio, error) {
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

// SetStsImgWithLabel ...
func SetStsImgWithLabel(id *ContainerPath, label ...string) error {
	return setStsImg(id, label...)
}

// SetStsImg update statefulset image field
func SetStsImg(id *ContainerPath) error {
	return setStsImg(id)
}

func setStsImg(id *ContainerPath, labels ...string) error {
	ctx := context.Background()
	opts := v1.GetOptions{}
	s, err := classicalClientSet.AppsV1().StatefulSets(id.Ns).Get(ctx, id.Name, opts)
	if err != nil {
		jsonlog.Err(err, map[string]any{"sts": id.Name, "msg": "get sts failed"})
		return err
	}
	if labels != nil {
		ls := s.GetLabels()
		if err = checkLabels(ls, labels); err != nil {
			return err
		}
	}
	cpy := s.DeepCopy()
	found := false
	for i, c := range cpy.Spec.Template.Spec.Containers {
		if c.Name == id.CName {
			jsonlog.Info("got new image", map[string]any{"sts": id.Name, "newImg": id.Img})
			cpy.Spec.Template.Spec.Containers[i].Image = id.Img
			found = true
			break
		}
	}
	if !found {
		fErr := errors.New("cannot find container")
		jsonlog.Err(fErr, map[string]any{"sts": id.Name, "container": id.CName, "msg": "cannot find container"})
		return fErr
	}
	uOpts := v1.UpdateOptions{}
	zero := int32(0)
	cpy.Spec.Replicas = &zero
	_, err = classicalClientSet.AppsV1().StatefulSets(id.Ns).Update(ctx, cpy, uOpts)
	if err != nil {
		jsonlog.Err(err, map[string]any{"sts": id.Name, "msg": "update sts failed"})
		return err
	}
	go func() {
		time.Sleep(5 * time.Second)
		a, rerr := classicalClientSet.AppsV1().StatefulSets(id.Ns).Get(ctx, id.Name, opts)
		if rerr != nil {
			jsonlog.Err(rerr, map[string]any{"name": id.Name, "msg": "get sts failed"})
			return
		}
		acpy := a.DeepCopy()
		acpy.Spec.Replicas = s.Spec.Replicas
		jsonlog.Debug("going to scale sts back replics",
			map[string]any{"*replicas": *s.Spec.Replicas, "replicas": s.Spec.Replicas})
		_, rerr = classicalClientSet.AppsV1().StatefulSets(id.Ns).Update(ctx, acpy, uOpts)
		if rerr != nil {
			jsonlog.Err(rerr, map[string]any{"name": id.Name, "msg": "scale sts failed"})
			return
		}
	}()
	return nil
}
