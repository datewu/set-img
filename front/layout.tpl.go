package front

import (
	"html/template"

	_ "embed"
)

//go:embed layout.html
var layoutHtml string

var layoutTpl = template.Must(template.New("layout").Parse(layoutHtml))
