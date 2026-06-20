package front

import (
	"html/template"
	"io"
	"path/filepath"
)

var indexTpl *template.Template
var indexTplWithLayout *template.Template

func initIndex() error {
	indexTpl = template.Must(template.New("index").ParseFiles(filepath.Join(rootDir, "index.html")))
	t, err := indexTpl.Clone()
	if err != nil {
		return err
	}
	indexTplWithLayout, err = t.AddParseTree("full index with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return nil
}

// IndexView ...
type IndexView struct {
	User  string
	Sites []IngressSite
}

// IngressSite mirrors k8s.IngressSite for template use.
type IngressSite struct {
	Ns        string
	Subdomain string
}

// Namespaces returns the deduplicated list of namespaces from Sites.
func (i IndexView) Namespaces() []string {
	seen := make(map[string]struct{})
	var ns []string
	for _, s := range i.Sites {
		if _, exists := seen[s.Ns]; !exists {
			seen[s.Ns] = struct{}{}
			ns = append(ns, s.Ns)
		}
	}
	return ns
}

// Render ...
func (i IndexView) Render(w io.Writer) error {
	return indexTpl.ExecuteTemplate(w, "content", i)
}

// FullPageRender ...
func (i IndexView) FullPageRender(w io.Writer, l LayoutView) error {
	l.Content = i
	return l.render(w, indexTplWithLayout)
}
