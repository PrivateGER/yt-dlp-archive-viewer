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
	Thumbnail	string
	Id 			string
	Extension 	string
	Metadata 	DirectoryIndexers.Metadata
}

func View(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList, path string) {
	FL.Lock()
	defer FL.Unlock()

	keys, ok := request.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	if _, ok := FL.Files[keys[0]]; !ok {
		return
	}

	// if no metadata loaded, do so
	if FL.Files[keys[0]].Metadata.ChannelID == "" {
		metadata, err := DirectoryIndexers.LoadMetadata(FL.Files[keys[0]], path)
		if err == nil {
			var fileObject DirectoryIndexers.VideoFile
			fileObject = FL.Files[keys[0]]
			fileObject.Metadata = metadata
			FL.Files[keys[0]] = fileObject
		}
	}

	video := FL.Files[keys[0]]
	data := ViewPageData{
		Title:     video.Title,
		Filename:  video.Filename,
		Thumbnail: video.Thumbnail,
		Id:        video.Id,
		Extension: video.Extension,
		Metadata:  video.Metadata,
	}

	err := tmpl["view.html"].ExecuteTemplate(writer, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}