package PageHandlers

import (
	_ "embed"
	"fmt"
	"net/http"
	"strconv"
	"ytdlp-viewer/DirectoryIndexers"
)

type IndexPageData struct {
	FileCount  string
	Files      map[string]DirectoryIndexers.VideoFile
	ShowingAll bool
}

func Index(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList) {
	FL.RLock()
	defer FL.RUnlock()

	keys, _ := request.URL.Query()["all"]

	if len(keys) == 0 {
		keys = []string{"0"}
	}

	var data IndexPageData
	if keys[0] == "1" {
		data = IndexPageData{
			FileCount:  strconv.Itoa(len(FL.Files)),
			Files:      FL.Files,
			ShowingAll: true,
		}
	} else {
		cutFiles := make(map[string]DirectoryIndexers.VideoFile)
		counter := 0

		for _, file := range FL.Files {
			if counter < 100 {
				cutFiles[file.Filename] = file
			}
			counter++
		}

		data = IndexPageData{
			FileCount:  strconv.Itoa(len(FL.Files)),
			Files:      cutFiles,
			ShowingAll: false,
		}
	}

	err := tmpl["index.html"].ExecuteTemplate(writer, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}
