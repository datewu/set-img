package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index.html
var indexHtml string

var indexTpl = template.Must(template.New("index").Parse(indexHtml))
var indexTplWithLayout *template.Template

func init() {
	t, err := indexTpl.Clone()
	if err != nil {
		panic(err)
	}
	indexTplWithLayout, err = t.AddParseTree("full index with layout", layoutTpl.Tree.Copy())
	if err != nil {
		panic(err)
	}
}

// IndexView ...
type IndexView struct {
	User string
}

func (i IndexView) Render(w io.Writer) error {
	return indexTplWithLayout.Execute(w, i)
}
