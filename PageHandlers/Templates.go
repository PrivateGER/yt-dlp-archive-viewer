package PageHandlers

import (
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed templates/index.html
var indexTmplSource string

//go:embed templates/view.html
var viewTmplSource string

//go:embed templates/search.html
var searchTmplSource string

//go:embed templates/base.html
var baseTmplSource string

var tmpl map[string]*template.Template

func init() {
	tmpl = make(map[string]*template.Template)

	var err error
	tmpl["index.html"] = template.New("index.html")
	tmpl["index.html"], err = tmpl["index.html"].Parse(baseTmplSource)
	if err != nil {
		fmt.Println(err)
	}
	tmpl["index.html"], err = tmpl["index.html"].Parse(indexTmplSource)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(tmpl["index.html"].Name())

	tmpl["search.html"] = template.New("search.html")
	tmpl["search.html"], err = tmpl["search.html"].Parse(searchTmplSource)
	if err != nil {
		fmt.Println(err)
	}
	tmpl["search.html"], err = tmpl["search.html"].Parse(baseTmplSource)
	if err != nil {
		fmt.Println(err)
	}

	tmpl["view.html"] = template.New("view.html")
	tmpl["view.html"], err = tmpl["view.html"].Parse(viewTmplSource)
	if err != nil {
		fmt.Println(err)
	}
	tmpl["view.html"], err = tmpl["view.html"].Parse(baseTmplSource)
	if err != nil {
		fmt.Println(err)
	}
}
