package PageHandlers

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"ytdlp-viewer/DirectoryIndexers"
)

type IndexPageData struct {
	FileCount 	string
	Files		map[string]DirectoryIndexers.VideoFile
}

func Index(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList) {
	FL.RLock()
	defer FL.RUnlock()

	data := IndexPageData{
		FileCount: strconv.Itoa(len(FL.Files)),
		Files:     FL.Files,
	}

	fmt.Println(tmpl["index.html"].Name())

	err := tmpl["index.html"].ExecuteTemplate(writer, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}
