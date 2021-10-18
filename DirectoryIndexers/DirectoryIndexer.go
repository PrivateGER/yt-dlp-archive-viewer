package DirectoryIndexers

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileList struct {
	Files map[string]VideoFile
	*sync.RWMutex
}

type VideoFile struct {
	Filename string
	Thumbnail string
	Extension string
	Title string
	Id string
	Metadata Metadata
}

func Index(path string, results chan FileList, oldFileList *FileList) {
	var FL FileList

	// Initialize the RWMutex HERE manually because *IT IS A POINTER TO A MUTEX*, so it defaults to a nil value
	FL.RWMutex = &sync.RWMutex{}

	fmt.Println("Scanning archive...")

	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return
	}

	FL.Files = make(map[string]VideoFile)
	var wg sync.WaitGroup

	bar := progressbar.NewOptions(len(fileList),
		progressbar.OptionSetDescription("Scanning files + metadata..."),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts())

	for _, video := range fileList {
		wg.Add(1)
		go func(video os.FileInfo) {
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
				wg.Done()
				_ = bar.Add(1)
				return
			}

			id := filenameToID(video.Name())

			// file has already been added before
			if _, ok := oldFileList.Files[id]; ok {
				FL.Lock()
				oldFileList.RLock()
				FL.Files[id] = oldFileList.Files[id]
				FL.Unlock()
				oldFileList.RUnlock()

				wg.Done()
				return
			}

			videoObject := VideoFile{
				Filename:  video.Name(),
				Thumbnail: strings.TrimSuffix(video.Name(), filepath.Ext(video.Name())) + ".webp",
				Extension: extension,
				Title:     filenameToTitle(video.Name(), extension),
				Id:        id,
				Metadata:  Metadata{},
			}

			metadata, err := LoadMetadata(videoObject, path)
			if err == nil {
				videoObject.Metadata = Metadata{
					Channel:       metadata.Channel,
					Thumbnail:     metadata.Thumbnail,
				}
			}

			FL.Lock()
			FL.Files[id] = videoObject
			FL.Unlock()

			_ = bar.Add(1)
			wg.Done()
		}(video)
	}

	wg.Wait()
	if !bar.IsFinished() {
		_ = bar.Finish()
	}

	results <- FL
	close(results)

	fmt.Println("\nArchive scan finished.")
}

func filenameToID(filename string) string {
	r := regexp2.MustCompile("([^[]+(?=]))", regexp2.RegexOptions(0))

	matches := regexp2FindAllString(r, filename)
	if len(matches) == 0 {
		r = regexp2.MustCompile("-[A-Za-z0-9_-]{11}", regexp2.RegexOptions(0))
		matches = regexp2FindAllString(r, filename)
		if len(matches) == 0 {
			log.Fatal("Got a video file without ID. Please remove or fix: " + filename)
		}

		// strips first dash away from the result (yes, I know this is DIRTY.
		matches = []string{matches[len(matches)-1]}
		matches[0] = matches[0][1:]
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