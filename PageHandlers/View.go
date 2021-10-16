package PageHandlers

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"ytdlp-viewer/DirectoryIndexers"
)

//go:embed templates/view.html
var viewTmplSource string
var viewTmpl *template.Template

type ViewPageData struct {
	Title 		string
	Filename 	string
	Id 			string
	Extension 	string
}

func init() {
	var err error
	viewTmpl = template.New("view.tmpl")
	viewTmpl, err = viewTmpl.Parse(viewTmplSource)
	if err != nil {
		fmt.Println(err)
	}
}

func View(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList) {
	FL.RLock()
	defer FL.RUnlock()

	keys, ok := request.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	if _, ok := FL.Files[keys[0]]; !ok {
		return
	}

	video := FL.Files[keys[0]]
	data := ViewPageData{
		Title:     video.Title,
		Filename:  video.Filename,
		Id:        video.Id,
		Extension: video.Extension,
	}

	err := viewTmpl.Execute(writer, data)
	if err != nil {
		fmt.Println(err)
	}
}