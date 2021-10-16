package PageHandlers

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"ytdlp-viewer/DirectoryIndexers"
)

type ViewPageData struct {
	Title 		string
	Filename 	string
	Id 			string
	Extension 	string
	Metadata 	DirectoryIndexers.Metadata
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
		Metadata:  video.Metadata,
	}

	err := tmpl["view.html"].ExecuteTemplate(writer, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}