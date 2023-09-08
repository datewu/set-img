package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index-layout.html
var indexHtml string

var indexTpl = template.Must(template.New("index").
	Delims("{i{", "}i}").Parse(indexHtml))

// IndexView ...
type IndexView struct {
	User string
}

func (i IndexView) Render(w io.Writer) {
	indexTpl.Execute(w, i)
}
