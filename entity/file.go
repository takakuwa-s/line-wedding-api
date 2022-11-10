package entity

import (
	"time"
)

type FileType string

const (
	Image = FileType("image")
	Video = FileType("video")
)

type FileStatus string

const (
	New      = FileStatus("new")
	Uploaded = FileStatus("uploaded")
	Open     = FileStatus("open")
	Deleted  = FileStatus("deleted")
)

type File struct {
	Id                 string     `json:"id"`
	Name               string     `json:"name"`
	FileType           FileType   `json:"fileType"`
	ContentUrl         string     `json:"contentUrl"`
	ThumbnailUrl       string     `json:"thumbnailUrl"`
	Width              int        `json:"width"`
	Height             int        `json:"height"`
	MimeType           string     `json:"mimeType"`
	FileStatus         FileStatus `json:"fileStatus"`
	Duration           int        `json:"duration"`
	FaceCount          int        `json:"faceCount"`
	FaceHappinessLevel float32    `json:"faceHappinessLevel"`
	FacePhotoBeauty    float32    `json:"facePhotoBeauty"`
	FaceScore          float32    `json:"faceScore"`
	ForBrideAndGroom   bool       `json:"forBrideAndGroom"`
	CreatedAt          time.Time  `json:"createdAt"`
	Creater            string     `json:"creater"`
	UpdatedAt          time.Time  `json:"updatedAt"`
}

func NewFile(messageId, creater string, fileType FileType, duration int) *File {
	return &File{
		Id:               messageId,
		FileType:         fileType,
		Duration:         duration,
		FileStatus:       New,
		ForBrideAndGroom: false,
		CreatedAt:        time.Now(),
		Creater:          creater,
		UpdatedAt:        time.Now(),
	}
}
