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
