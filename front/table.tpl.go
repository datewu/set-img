package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed table.html
var tableHtml string

var tableTpl = template.Must(template.New("table").Parse(tableHtml))

// TableView ...
type TableView struct {
	TODO string
}

func (t TableView) Render(w io.Writer) {
	tableTpl.Execute(w, t)
}
