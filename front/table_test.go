package front

import (
	"bytes"
	"testing"
)

func TestTableViewRender(t *testing.T) {
	view := TableView{
		Description: "test description",
		Namespace:   "test-ns",
		Kind:        "deploy",
		Data: []Resource{
			{
				Name:      "test-deploy",
				Namespace: "test-ns",
				Replicas:  3,
				Age:       "10d",
				Containers: []Container{
					{
						Name:  "web",
						Image: "nginx:latest",
						Env: []EnvKeyVal{
							{Key: "ENV_VAR", Value: "val"},
						},
					},
					{
						Name:  "sidecar",
						Image: "busybox:latest",
						Env:   []EnvKeyVal{},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	err := view.Render(&buf)
	if err != nil {
		t.Fatalf("Failed to render table content: %v", err)
	}

	layout := NewLayout("test-user", "development")
	var bufFull bytes.Buffer
	err = view.FullPageRender(&bufFull, layout)
	if err != nil {
		t.Fatalf("Failed to render full page: %v", err)
	}

	t.Logf("Rendered full page length: %d", bufFull.Len())
}
