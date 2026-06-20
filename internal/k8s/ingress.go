package k8s

import (
	"context"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListIngressHostsByLabel lists ingresses matching the label selector
// and returns deduplicated subdomains (the part before .deoops.com).
func ListIngressHostsByLabel(ns, label string) ([]string, error) {
	ctx := context.Background()
	opts := v1.ListOptions{LabelSelector: label}
	ingresses, err := classicalClientSet.NetworkingV1().Ingresses(ns).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	var subdomains []string
	for _, ing := range ingresses.Items {
		for _, rule := range ing.Spec.Rules {
			host := rule.Host
			sub, ok := strings.CutSuffix(host, ".deoops.com")
			if !ok || sub == "" {
				continue
			}
			if _, exists := seen[sub]; !exists {
				seen[sub] = struct{}{}
				subdomains = append(subdomains, sub)
			}
		}
	}
	return subdomains, nil
}
