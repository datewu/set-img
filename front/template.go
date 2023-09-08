package front

import (
	"html/template"

	_ "embed"
)

//go:embed index-layout.html
var s string

// IndexTpl for index-layout.html
var IndexTpl = template.Must(template.New("index").
	Delims("{i{", "}i}").Parse(s))
