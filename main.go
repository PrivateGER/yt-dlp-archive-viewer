package main

import (
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"ytdlp-viewer/DirectoryIndexers"
	"ytdlp-viewer/PageHandlers"
)

func main() {
	path := os.Getenv("directory")
	if path == "" {
		path = "/home/latte/NFS/MainArchive/"
	}

	var FL DirectoryIndexers.FileList
	resultChannel := make(chan DirectoryIndexers.FileList)

	fmt.Println("Starting scanner at", path)
	go DirectoryIndexers.Index(path, resultChannel)

	FL = <-resultChannel

	/*FL.RLock()
	fmt.Println(strconv.Itoa(len(FL.Files)))
	for _, file := range FL.Files {
		fmt.Println("Name:", file.Title,"Extension:", file.Extension, "ID:", file.Id)
	}
	FL.RUnlock()*/

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		PageHandlers.Index(writer, request, &FL)
	})
	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		PageHandlers.SearchHandler(writer, request, &FL)
	})
	http.HandleFunc("/view", func(writer http.ResponseWriter, request *http.Request) {
		PageHandlers.View(writer, request, &FL)
	})
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(path))))

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
