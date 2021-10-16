package PageHandlers

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"ytdlp-viewer/DirectoryIndexers"
)

//go:embed templates/index.html
var indexTmplSource string
var indexTmpl *template.Template

type IndexPageData struct {
	FileCount 	string
	Files		map[string]DirectoryIndexers.VideoFile
}

func init() {
	var err error
	indexTmpl = template.New("index.tmpl")
	indexTmpl, err = indexTmpl.Parse(indexTmplSource)
	if err != nil {
		fmt.Println(err)
	}
}

func Index(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList) {
	FL.RLock()
	defer FL.RUnlock()

	data := IndexPageData{
		FileCount: strconv.Itoa(len(FL.Files)),
		Files:     FL.Files,
	}

	err := indexTmpl.Execute(writer, data)
	if err != nil {
		fmt.Println(err)
	}
}
