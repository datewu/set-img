package front

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"time"

	apps "k8s.io/api/apps/v1"
)

var tableTpl *template.Template

var tableTplWithLayout *template.Template

func initTable() error {
	tableTpl = template.Must(template.New("table").ParseFiles(filepath.Join(rootDir, "table.html")))
	t, err := tableTpl.Clone()
	if err != nil {
		return err
	}
	tableTplWithLayout, err = t.AddParseTree("full table with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return nil
}

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
	Age        string
}

type Container struct {
	Name, Image string
}

func (r *Resource) formatAge(t time.Time) {
	d := time.Now().Round(time.Hour).Sub(t.Round(time.Hour))
	age := ""
	if d.Hours() > 24 {
		days := d.Hours() / 24
		if days > 365 {
			years := days / 365
			y := int(years)
			age = fmt.Sprintf("%dy", y)
			d -= time.Duration(y*365*24) * time.Hour
			days = d.Hours() / 24
		}
		if days > 1 {
			i := int(days)
			age = fmt.Sprintf("%dd", i)
			d -= time.Duration(i*24) * time.Hour
		}
	}
	age += d.String()
	r.Age = strings.ReplaceAll(age, "0m0s", "")
}

func newDeployResource(d *apps.Deployment) *Resource {
	res := &Resource{
		Name:     d.Name,
		Replicas: int(*d.Spec.Replicas),
	}
	res.formatAge(d.ObjectMeta.GetCreationTimestamp().Time)
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
	}
	res.formatAge(s.ObjectMeta.GetCreationTimestamp().Time)
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

func (t TableView) Render(w io.Writer) error {
	return tableTpl.ExecuteTemplate(w, "content", t)
}

// FullPageRender ...
func (t TableView) FullPageRender(w io.Writer, l LayoutView) error {
	l.Content = t
	return l.render(w, tableTplWithLayout)
}
