package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed layout.html
var layoutHtml string

var layoutTpl = template.Must(template.New("layout").Parse(layoutHtml))

// LayoutView is a view for the layou
type LayoutView struct {
	User    string
	Env     string
	Content any
}

func NewLayout(user, env string) LayoutView {
	return LayoutView{
		User: user,
		Env:  env,
	}
}

func (l LayoutView) render(w io.Writer, tpl *template.Template) error {
	return tpl.Execute(w, l)
}

// InitOrReload init or reload the layout
func InitOrReload() error {
	if err := initIndex(); err != nil {
		return err
	}

	if err := initProfile(); err != nil {
		return err
	}

	if err := initTable(); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := InitOrReload(); err != nil {
		panic(err)
	}
}
