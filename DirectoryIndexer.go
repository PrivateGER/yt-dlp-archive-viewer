package main

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

type FileList struct {
	files map[string]VideoFile
	mu sync.Mutex
}

type VideoFile struct {
	Filename string
	Extension string
	Title string
	Id string
}

var FL FileList

func ScanDirectory(group *sync.WaitGroup, path string) {
	defer group.Done()

	FL.mu.Lock()
	defer FL.mu.Unlock()

	fmt.Println("Scanning archive...")

	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return
	}

	FL.files = make(map[string]VideoFile)
	for _, video := range fileList {
		extension := filepath.Ext(video.Name())[1:]
		// check if extension is one of valid yt-dlp extensions, if not ignore file
		switch extension {
		case "3gp",
			"aac",
			"flv",
			"m4a",
			"mp3",
			"mp4",
			"webm":
				break
		default:
			continue
		}

		id := filenameToID(video.Name())

		FL.files[id] = VideoFile{
			Filename:  video.Name(),
			Extension: extension,
			Title:     filenameToTitle(video.Name(), extension),
			Id:        id,
		}
	}

	fmt.Println("Archive scan finished.")
}

func filenameToID(filename string) string {
	r := regexp2.MustCompile("([^[]+(?=]))", regexp2.RegexOptions(0))

	matches := regexp2FindAllString(r, filename)
	if len(matches) == 0 {
		fmt.Println("Got video without square-bracket ID format. Falling back to youtube-dl 11-char string matching (!THIS MAY CAUSE ISSUES!):", filename)

		r = regexp2.MustCompile("-[A-Za-z0-9_-]{11}", regexp2.RegexOptions(0))
		matches = regexp2FindAllString(r, filename)
		if len(matches) == 0 {
			log.Fatal("Got a video file without ID. Please remove or fix: " + filename)
		}

		// strips first dash away from the result (yes, I know this is DIRTY.
		matches = []string{matches[len(matches)-1]}
		matches[0] = matches[0][1:]
		fmt.Println("Recovered ID:", matches[0])
	}
	return matches[len(matches)-1] // last element = the id between square brackets
}

func filenameToTitle(filename string, extension string) string {
	title := strings.Replace(filename, "[" + filenameToID(filename) + "]", "", 1)
	title = strings.Replace(title, "." + extension, "", 1)
	title = strings.Replace(title, filenameToID(filename), "", 1)

	return title
}

func regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}