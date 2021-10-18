package main

import (
	_ "embed"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
	"ytdlp-viewer/DirectoryIndexers"
	"ytdlp-viewer/PageHandlers"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "./", "the full path to the ytdlp archive (with a / suffix)")
	var port int
	flag.IntVar(&port, "port", 8000, "the port for the web panel to listen on")

	flag.Parse()

	var FL DirectoryIndexers.FileList
	FL.RWMutex = &sync.RWMutex{}

	go func() {
		for {
			resultChannel := make(chan DirectoryIndexers.FileList)

			fmt.Println("Starting scanner at", path)
			go DirectoryIndexers.Index(path, resultChannel, &FL)

			refreshedFL := <-resultChannel
			FL.Lock()
			FL.Files = refreshedFL.Files
			FL.Unlock()
			fmt.Println("File list refreshed")

			time.Sleep(60 * time.Second)
		}
	}()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(),request.URL)
		PageHandlers.Index(writer, request, &FL)
	})
	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(),request.URL)
		PageHandlers.SearchHandler(writer, request, &FL)
	})
	http.HandleFunc("/view", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(),request.URL)
		PageHandlers.View(writer, request, &FL, path)
	})
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(path))))

	err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
