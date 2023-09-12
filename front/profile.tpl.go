package front

import (
	"html/template"
	"io"

	_ "embed"
)

var profileHtml = `
{{- define "content"  }}
  <span>welcome {{ .User }}</span>
{{ end -}}
`

// profileTpl is the profile template.
var profileTpl = template.Must(template.New("profile").Parse(profileHtml))

// ProfileView ...
type ProfileView struct {
	User string
}

func (p ProfileView) Render(w io.Writer) error {
	t, err := profileTpl.Clone()
	if err != nil {
		return err
	}
	t, err = t.AddParseTree("full index with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return t.Execute(w, p)
}
