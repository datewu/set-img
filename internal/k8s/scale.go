package k8s

import (
	"context"
	"errors"

	"github.com/datewu/gtea/jsonlog"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScalePath identifies a deployment/statefulset and the desired replica count.
type ScalePath struct {
	Ns       string
	Kind     string
	Name     string
	Replicas int32
}

func (s ScalePath) valid() bool {
	if s.Ns == "" || s.Name == "" || s.Replicas < 0 {
		return false
	}
	if s.Kind != "deploy" && s.Kind != "sts" {
		return false
	}
	return true
}

// Scale updates the replica count of a deployment or statefulset.
func (s ScalePath) Scale(labels ...string) error {
	if !s.valid() {
		return errors.New("invalid scale param")
	}
	ctx := context.Background()
	opts := v1.GetOptions{}
	uOpts := v1.UpdateOptions{}
	switch s.Kind {
	case "deploy":
		d, err := classicalClientSet.AppsV1().Deployments(s.Ns).Get(ctx, s.Name, opts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"deploy": s.Name, "msg": "get deploy failed"})
			return err
		}
		if err = checkLabels(d.GetLabels(), labels); err != nil {
			return err
		}
		sc, err := classicalClientSet.AppsV1().Deployments(s.Ns).GetScale(ctx, s.Name, opts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"deploy": s.Name, "msg": "get deploy scale failed"})
			return err
		}
		sc.Spec.Replicas = s.Replicas
		_, err = classicalClientSet.AppsV1().Deployments(s.Ns).UpdateScale(ctx, s.Name, sc, uOpts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"deploy": s.Name, "replicas": s.Replicas, "msg": "scale deploy failed"})
			return err
		}
	case "sts":
		st, err := classicalClientSet.AppsV1().StatefulSets(s.Ns).Get(ctx, s.Name, opts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"sts": s.Name, "msg": "get sts failed"})
			return err
		}
		if err = checkLabels(st.GetLabels(), labels); err != nil {
			return err
		}
		sc, err := classicalClientSet.AppsV1().StatefulSets(s.Ns).GetScale(ctx, s.Name, opts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"sts": s.Name, "msg": "get sts scale failed"})
			return err
		}
		sc.Spec.Replicas = s.Replicas
		_, err = classicalClientSet.AppsV1().StatefulSets(s.Ns).UpdateScale(ctx, s.Name, sc, uOpts)
		if err != nil {
			jsonlog.Err(err, map[string]any{"sts": s.Name, "replicas": s.Replicas, "msg": "scale sts failed"})
			return err
		}
	}
	return nil
}
