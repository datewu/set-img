package k8s

import "testing"

func TestScalePathInvalid(t *testing.T) {
	cases := []struct {
		name string
		path ScalePath
	}{
		{"empty namespace", ScalePath{Kind: "deploy", Name: "web", Replicas: 1}},
		{"empty name", ScalePath{Ns: "wu", Kind: "deploy", Replicas: 1}},
		{"unknown kind", ScalePath{Ns: "wu", Kind: "pod", Name: "web", Replicas: 1}},
		{"negative replicas", ScalePath{Ns: "wu", Kind: "deploy", Name: "web", Replicas: -1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.path.Scale(); err == nil {
				t.Fatalf("expected an error for %+v", c.path)
			}
		})
	}
}
