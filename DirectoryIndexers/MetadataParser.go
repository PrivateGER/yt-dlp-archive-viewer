package DirectoryIndexers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Metadata struct {
	Description   string     `json:"description"`
	UploadDate    string     `json:"upload_date"`
	UploaderURL   string     `json:"uploader_url"`
	ChannelID     string     `json:"channel_id"`
	ChannelURL    string     `json:"channel_url"`
	ViewCount     int        `json:"view_count"`
	AverageRating float64    `json:"average_rating"`
	AgeLimit      int        `json:"age_limit"`
	WebpageURL    string     `json:"webpage_url"`
	LikeCount     int        `json:"like_count"`
	DislikeCount  int        `json:"dislike_count"`
	Channel       string     `json:"channel"`
	Thumbnail     string     `json:"thumbnail"`
	DisplayID     string     `json:"display_id"`
	Width         int        `json:"width"`
	Height        int        `json:"height"`
	Fps           int        `json:"fps"`
	Vcodec        string     `json:"vcodec"`
	Acodec        string     `json:"acodec"`
	Abr           float64    `json:"abr"`
	Fulltitle     string     `json:"fulltitle"`
	Comments      []Comments `json:"comments"`
	CommentCount  int        `json:"comment_count"`
}
type Comments struct {
	ID               string `json:"id"`
	Text             string `json:"text"`
	Timestamp        int    `json:"timestamp"`
	LikeCount        int    `json:"like_count"`
	Author           string `json:"author"`
	AuthorID         string `json:"author_id"`
	AuthorThumbnail  string `json:"author_thumbnail"`
	AuthorIsUploader bool   `json:"author_is_uploader"`
}

func ParseMetadata(jsonBytes []byte) (Metadata, error) {
	var meta Metadata
	err := json.Unmarshal(jsonBytes, &meta)
	if err != nil {
		fmt.Println(err)
		return Metadata{}, err
	}
	return meta, err
}

func LoadMetadata(videoobject VideoFile, path string) (Metadata, error) {
	metadataFilename := strings.TrimSuffix(videoobject.Filename, filepath.Ext(videoobject.Filename)) + ".info.json"
	jsonMetadataFile, err := os.Open(path + metadataFilename)
	if err != nil {
		return Metadata{}, err
	}

	metadataBytes, _ := ioutil.ReadAll(jsonMetadataFile)
	_ = jsonMetadataFile.Close()

	metadata, err := ParseMetadata(metadataBytes)

	// cut off file extension to get the webp thumbnail filename
	thumbnailFilename := strings.TrimSuffix(path + videoobject.Filename, filepath.Ext(videoobject.Filename)) + ".webp"
	if _, err := os.Stat(thumbnailFilename); errors.Is(err, os.ErrNotExist) {
		metadata.Thumbnail = ""
	} else {
		metadata.Thumbnail = strings.TrimSuffix(videoobject.Filename, filepath.Ext(videoobject.Filename)) + ".webp"
	}

	if err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}
