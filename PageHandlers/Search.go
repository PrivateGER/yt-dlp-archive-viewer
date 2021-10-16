package PageHandlers

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"ytdlp-viewer/DirectoryIndexers"
)

type SearchPageData struct {
	Results		[]DirectoryIndexers.VideoFile
	ResultCount	string
	SearchTerm	string
}

func SearchHandler(writer http.ResponseWriter, request *http.Request, FL *DirectoryIndexers.FileList) {
	FL.RLock()
	defer FL.RUnlock()

	keys, ok := request.URL.Query()["term"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'term' is missing")
		return
	}

	var results []DirectoryIndexers.VideoFile

	for _, video := range FL.Files {
		if video.Id == keys[0] {
			results = append(results, video)
			break
		}
		if strings.Contains(strings.ToUpper(video.Title), strings.ToUpper(keys[0])) {
			results = append(results, video)
			continue
		}
	}

	data := SearchPageData{
		Results:     results,
		ResultCount: strconv.Itoa(len(results)),
		SearchTerm: keys[0],
	}

	err := tmpl["search.html"].ExecuteTemplate(writer, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}
