package k8s

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressSite represents a subdomain from an ingress in a specific namespace.
type IngressSite struct {
	Ns        string
	Subdomain string
}

// ListIngressSitesByLabel lists ingresses matching the label selector
// and returns namespace + subdomain pairs (the part before .deoops.com).
func ListIngressSitesByLabel(ns, label string) ([]IngressSite, error) {
	ctx := context.Background()
	opts := v1.ListOptions{LabelSelector: label}
	ingresses, err := classicalClientSet.NetworkingV1().Ingresses(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	var sites []IngressSite
	for _, ing := range ingresses.Items {
		if len(ing.Spec.Rules) == 0 {
			continue
		}
		host := ing.Spec.Rules[0].Host
		sub, ok := strings.CutSuffix(host, ".deoops.com")
		if !ok || sub == "" {
			continue
		}
		key := ing.Namespace + "/" + sub
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			sites = append(sites, IngressSite{Ns: ing.Namespace, Subdomain: sub})
		}
	}
	return sites, nil
}
