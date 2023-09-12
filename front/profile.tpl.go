package front

import (
	"html/template"

	_ "embed"
)

var profileHtml = `
{{- define "content"  }}
  <span>welcome {{ .User }}</span>
{{ end -}}
`

// ProfileTpl is the profile template.
var ProfileTpl = template.Must(template.New("profile").Parse(profileHtml + layoutHtml))
