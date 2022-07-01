package entity

import (
	"time"
)

type FileType string

const (
	Image = FileType("image")
	Video = FileType("video")
)

type File struct {
	Id                 string    `json:"id"`
	Name               string    `json:"name"`
	FileType           FileType  `json:"fileType"`
	ContentUrl         string    `json:"contentUrl"`
	ThumbnailUrl       string    `json:"thumbnailUrl"`
	Width              int       `json:"width"`
	Height             int       `json:"height"`
	MimeType           string    `json:"mimeType"`
	Uploaded           bool      `json:"uploaded"`
	Duration           int       `json:"duration"`
	Calculated         bool      `json:"calculated"`
	FaceCount          int       `json:"faceCount"`
	FaceHappinessLevel float32   `json:"faceHappinessLevel"`
	FacePhotoBeauty    float32   `json:"facePhotoBeauty"`
	FaceScore          float32   `json:"faceScore"`
	CreatedAt          time.Time `json:"createdAt"`
	Creater            string    `json:"creater"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func NewFile(messageId, creater string, fileType FileType, duration int) *File {
	return &File{
		Id:         messageId,
		FileType:   fileType,
		Duration:   duration,
		Uploaded:   false,
		Calculated: false,
		CreatedAt:  time.Now(),
		Creater:    creater,
		UpdatedAt:  time.Now(),
	}
}
