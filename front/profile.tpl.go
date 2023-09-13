package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed profile.html
var profileHtml string

// profileTpl is the profile template.
var profileTpl = template.Must(template.New("profile").Parse(profileHtml))
var profileTplWithLayout *template.Template

func init() {
	t, err := profileTpl.Clone()
	if err != nil {
		panic(err)
	}
	profileTplWithLayout, err = t.AddParseTree("full profile with layout", layoutTpl.Tree.Copy())
	if err != nil {
		panic(err)
	}
}

// ProfileView ...
type ProfileView struct {
	User string
}

func (p ProfileView) Render(w io.Writer) error {
	return profileTplWithLayout.Execute(w, p)
}
