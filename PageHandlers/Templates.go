package PageHandlers

import (
	"embed"
	_ "embed"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

var tmpl map[string]*template.Template

func init() {
	tmpl = make(map[string]*template.Template)

	tmpl["index.html"] = template.Must(template.ParseFS(templateFS, "templates/index.html", "templates/base.html"))
	tmpl["search.html"] = template.Must(template.ParseFS(templateFS, "templates/search.html", "templates/base.html"))
	tmpl["view.html"] = template.Must(template.ParseFS(templateFS, "templates/view.html", "templates/base.html"))
}
