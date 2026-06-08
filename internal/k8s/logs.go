package k8s

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// PodContainer identifies a specific container within the pods of a workload
type PodContainer struct {
	Ns        string
	Kind      string // "deploy" or "sts"
	Name      string // deployment/statefulset name
	Container string // container name
}

// listPodsByOwner lists pods that belong to a deployment or statefulset
func listPodsByOwner(ctx context.Context, ns, kind, name string) ([]corev1.Pod, error) {
	var labelSelector string
	switch kind {
	case "deploy":
		d, err := classicalClientSet.AppsV1().Deployments(ns).Get(ctx, name, v1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get deployment %s/%s: %w", ns, name, err)
		}
		labelSelector = labels.Set(d.Spec.Selector.MatchLabels).String()
	case "sts":
		s, err := classicalClientSet.AppsV1().StatefulSets(ns).Get(ctx, name, v1.GetOptions{})
		if err != nil {
			return nil, fmt.Errorf("get statefulset %s/%s: %w", ns, name, err)
		}
		labelSelector = labels.Set(s.Spec.Selector.MatchLabels).String()
	default:
		return nil, fmt.Errorf("unsupported kind: %s", kind)
	}

	pods, err := classicalClientSet.CoreV1().Pods(ns).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, fmt.Errorf("list pods for %s/%s: %w", kind, name, err)
	}
	return pods.Items, nil
}

// StreamPodLogs returns an io.ReadCloser that streams logs from a specific
// pod and container. The caller is responsible for closing the reader.
func StreamPodLogs(ctx context.Context, ns, podName, container string, tailLines int64) (io.ReadCloser, error) {
	opts := &corev1.PodLogOptions{
		Container: container,
		Follow:    true,
		TailLines: &tailLines,
	}
	req := classicalClientSet.CoreV1().Pods(ns).GetLogs(podName, opts)
	return req.Stream(ctx)
}

// ListWorkloadPods returns a list of pod names for a given workload (deploy/sts)
func ListWorkloadPods(ctx context.Context, ns, kind, name string) ([]string, error) {
	pods, err := listPodsByOwner(ctx, ns, kind, name)
	if err != nil {
		return nil, err
	}
	names := make([]string, len(pods))
	for i, p := range pods {
		names[i] = p.Name
	}
	return names, nil
}
