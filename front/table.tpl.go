package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index-layout.html
var t string

var tableTpl = template.Must(template.New("table").Parse(t))

// TableView ...
type TableView struct {
	TODO string
}

func (t TableView) Render(w io.Writer) {
	tableTpl.Execute(w, t)
}
