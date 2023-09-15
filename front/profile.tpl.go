package front

import (
	"html/template"
	"io"
	"path/filepath"
)

// profileTpl is the profile template.
var profileTpl *template.Template
var profileTplWithLayout *template.Template

func initProfile() error {
	profileTpl = template.Must(template.New("profile").ParseFiles(filepath.Join(rootDir, "profile.html")))
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

// FullPageRender ...
func (p ProfileView) FullPageRender(w io.Writer, l LayoutView) error {
	l.Content = p
	return l.render(w, profileTplWithLayout)
}
