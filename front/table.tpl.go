package front

import (
	"html/template"
	"io"
	"time"

	apps "k8s.io/api/apps/v1"

	_ "embed"
)

//go:embed table.html
var tableHtml string

var tableTpl = template.Must(template.New("table").Parse(tableHtml))

// TableView ...
type TableView struct {
	Description string
	Namespace   string
	Kind        string
	Data        []Resource
}

// AddDeploys ...
func (t *TableView) AddDeploys(ds []apps.Deployment) {
	var res []Resource
	for _, d := range ds {
		res = append(res, *newDeployResource(&d))
	}
	t.Data = res
}

// AddSts ...
func (t *TableView) AddSts(ss []apps.StatefulSet) {
	var res []Resource
	for _, s := range ss {
		res = append(res, *newStsResource(&s))
	}
	t.Data = res
}

type Resource struct {
	Containers []Container
	Name       string
	Replicas   int
	Age        time.Duration
}

func newDeployResource(d *apps.Deployment) *Resource {
	res := &Resource{
		Name:     d.Name,
		Replicas: int(*d.Spec.Replicas),
		Age: time.Now().Round(time.Second).
			Sub((d.ObjectMeta.GetCreationTimestamp().Round(time.Second))),
	}
	containes := d.Spec.Template.Spec.Containers
	cs := make([]Container, len(containes))
	for i, c := range containes {
		cs[i] = Container{
			Name:  c.Name,
			Image: c.Image,
		}
	}
	res.Containers = cs
	return res
}
func newStsResource(s *apps.StatefulSet) *Resource {
	res := &Resource{
		Name:     s.Name,
		Replicas: int(*s.Spec.Replicas),
		Age: time.Now().Round(time.Second).
			Sub((s.ObjectMeta.GetCreationTimestamp().Round(time.Second))),
	}
	containes := s.Spec.Template.Spec.Containers
	cs := make([]Container, len(containes))
	for i, c := range containes {
		cs[i] = Container{
			Name:  c.Name,
			Image: c.Image,
		}
	}
	res.Containers = cs
	return res
}

type Container struct {
	Name, Image string
}

func (t TableView) Render(w io.Writer) {
	tableTpl.Execute(w, t)
}
