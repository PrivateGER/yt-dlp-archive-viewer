package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

//go:embed templates/index.html
var indexTmplSource string

//go:embed templates/search.html
var searchTmplSource string

//go:embed templates/view.html
var viewTmplSource string

type IndexPageData struct {
	FileCount 	string
	Files		map[string]VideoFile
}

type SearchPageData struct {
	Results		[]VideoFile
	ResultCount	string
	SearchTerm	string
}

type ViewPageData struct {
	Title 		string
	Filename 	string
	Id 			string
	Extension 	string
}

func main() {
	path := os.Getenv("directory")
	if path == "" {
		path = "Z:/MainArchive"
	}

	var wg sync.WaitGroup
	wg.Add(1)

	fmt.Println("Starting scanner at", path)
	go ScanDirectory(&wg, path)

	wg.Wait()

	FL.mu.Lock()
	fmt.Println(strconv.Itoa(len(FL.files)))
	for _, file := range FL.files {
		fmt.Println("Name:", file.Title,"Extension:", file.Extension, "ID:", file.Id)
	}
	FL.mu.Unlock()

	indexTmpl, err := template.New("index.tmpl").Parse(indexTmplSource)
	if err != nil {
		fmt.Println(err)
	}
	searchTmpl, err := template.New("search.tmpl").Parse(searchTmplSource)
	if err != nil {
		fmt.Println(err)
	}
	viewTmpl, err := template.New("view.tmpl").Parse(viewTmplSource)
	if err != nil {
		fmt.Println(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		FL.mu.Lock()
		defer FL.mu.Unlock()

		data := IndexPageData{
			FileCount: strconv.Itoa(len(FL.files)),
			Files:     FL.files,
		}

		err := indexTmpl.Execute(writer, data)
		if err != nil {
			fmt.Println(err)
		}
	})
	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		FL.mu.Lock()
		defer FL.mu.Unlock()

		keys, ok := request.URL.Query()["term"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'term' is missing")
			return
		}

		var results []VideoFile

		for _, video := range FL.files {
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

		err := searchTmpl.Execute(writer, data)
		if err != nil {
			fmt.Println(err)
		}
	})
	http.HandleFunc("/view", func(writer http.ResponseWriter, request *http.Request) {
		FL.mu.Lock()
		defer FL.mu.Unlock()

		keys, ok := request.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			log.Println("Url Param 'id' is missing")
			return
		}

		if _, ok := FL.files[keys[0]]; !ok {
			return
		}

		video := FL.files[keys[0]]
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
	})
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(path))))

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
