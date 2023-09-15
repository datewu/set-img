package front

import (
	"html/template"
	"io"

	_ "embed"
)

//go:embed index.html
var indexHtml string

var indexTpl = template.Must(template.New("index").Parse(indexHtml))
var indexTplWithLayout *template.Template

func initIndex() error {
	t, err := indexTpl.Clone()
	if err != nil {
		return err
	}
	indexTplWithLayout, err = t.AddParseTree("full index with layout", layoutTpl.Tree.Copy())
	if err != nil {
		return err
	}
	return nil
}

// IndexView ...
type IndexView struct {
}

// Render ...
func (i IndexView) Render(w io.Writer) error {
	return indexTpl.ExecuteTemplate(w, "content", i)
}

// FullPage embed index in layout template
func (i IndexView) FullPage(user, env string) LayoutView {
	l := LayoutView{
		User:       user,
		Env:        env,
		ContentTpl: indexTplWithLayout,
		Content:    i,
	}
	return l
}
