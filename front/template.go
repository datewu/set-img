package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index-layout.html
var s string

var indexTpl = template.Must(template.New("index").
	Delims("{i{", "}i}").ParseFiles(s))

// IndexView ...
type IndexView struct {
	User string
}

func (i IndexView) Render(w io.Writer) {
	indexTpl.Execute(w, i)
}
