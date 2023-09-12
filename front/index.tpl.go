package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index.html
var indexHtml string

var indexTpl = template.Must(template.New("index").Parse(indexHtml))

// IndexView ...
type IndexView struct {
	User string
}

func (i IndexView) Render(w io.Writer) error {
	t, err := indexTpl.Clone()
	if err != nil {
		return err
	}
	t, err = t.AddParseTree("full index with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return t.Execute(w, i)
}
