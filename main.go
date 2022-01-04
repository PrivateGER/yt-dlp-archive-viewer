package main

import (
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"ytdlp-viewer/DirectoryIndexers"
	"ytdlp-viewer/PageHandlers"
)

//go:embed PageHandlers/static
var staticEmbedFS embed.FS

func main() {
	var path string
	flag.StringVar(&path, "path", "/archive", "the full path to the ytdlp archive")
	var port int
	flag.IntVar(&port, "port", 8000, "the port for the web panel to listen on")
	var autorefresh bool
	flag.BoolVar(&autorefresh, "refresh", false, "whether should the index be updated every x seconds")
	var refreshInterval int
	flag.IntVar(&refreshInterval, "interval", 60, "the interval for the index to update in seconds")

	flag.Parse()

	// append last slash in case it's not provided
	if !strings.HasSuffix(path, "/") {
		path += "/"
		fmt.Println("Added missing trailing slash.")
	}

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

			if !autorefresh {
				return
			}

			time.Sleep(time.Duration(refreshInterval) * time.Second)
		}
	}()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(), request.URL)
		PageHandlers.Index(writer, request, &FL)
	})
	http.HandleFunc("/search", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(), request.URL)
		PageHandlers.SearchHandler(writer, request, &FL)
	})
	http.HandleFunc("/view", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(time.Now().String(), request.URL)
		PageHandlers.View(writer, request, &FL, path)
	})

	subFS, err := fs.Sub(staticEmbedFS, "PageHandlers/static")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(subFS))))

	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir(path))))

	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
