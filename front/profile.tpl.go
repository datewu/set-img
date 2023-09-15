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

func initProfile() error {
	t, err := profileTpl.Clone()
	if err != nil {
		return err
	}
	profileTplWithLayout, err = t.AddParseTree("full profile with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return nil
}

// ProfileView ...
type ProfileView struct {
	User string
}

func (p ProfileView) Render(w io.Writer) error {
	return profileTpl.ExecuteTemplate(w, "content", p)
}

// FullPage embed profile with layout template
func (p ProfileView) FullPage(user, env string) LayoutView {
	l := LayoutView{
		User:       user,
		Env:        env,
		ContentTpl: profileTplWithLayout,
		Content:    p,
	}
	return l
}
