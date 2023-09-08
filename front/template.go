package front

import "html/template"

// IndexTpl for index-layout.html
var IndexTpl = template.Must(template.New("index").
	Delims("{i{", "}i}").ParseFiles("index-layout.html"))
