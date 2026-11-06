package front

import (
	"bytes"
	"strings"
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

	rendered := buf.String()
	for _, want := range []string{
		`hx-put="/my/scale"`,
		`name="replicas" value="3"`,
		`class="form-input replicas-input"`,
		`data-current="3"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered table is missing %q", want)
		}
	}

	layout := NewLayout("test-user", "development")
	var bufFull bytes.Buffer
	err = view.FullPageRender(&bufFull, layout)
	if err != nil {
		t.Fatalf("Failed to render full page: %v", err)
	}

	t.Logf("Rendered full page length: %d", bufFull.Len())
}
