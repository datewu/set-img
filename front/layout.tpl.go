package front

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
)

var rootDir = findRootDir()

func findRootDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "front"
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "front", "layout.html")); err == nil {
			return filepath.Join(dir, "front")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "front"
}

var layoutTpl *template.Template

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
	layout, err := os.ReadFile(filepath.Join(rootDir, "layout.html"))
	if err != nil {
		return err
	}
	layoutTpl = template.Must(template.New("layout").Parse(string(layout)))
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
