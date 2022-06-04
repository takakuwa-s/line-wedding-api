package entity

import (
	"time"
)

type File struct {
	Id                 string    `json:"id"`
	Name               string    `json:"name"`
	ContentUrl         string    `json:"contentUrl"`
	ThumbnailUrl       string    `json:"thumbnailUrl"`
	Width              int       `json:"width"`
	Height             int       `json:"height"`
	MimeType           string    `json:"mimeType"`
	IsUploaded         bool      `json:"isUploaded"`
	FaceCount          int       `json:"faceCount"`
	FaceHappinessLevel float32   `json:"faceHappinessLevel"`
	FacePhotoBeauty    float32   `json:"facePhotoBeauty"`
	FaceScore          float32   `json:"faceScore"`
	CreatedAt          time.Time `json:"createdAt"`
	Creater            string    `json:"creater"`
	UpdatedAt          time.Time `json:"updatedAt"`
	Updater            string    `json:"updater"`
}

func NewFile(messageId, creater string) *File {
	return &File{
		Id:         messageId,
		IsUploaded: false,
		CreatedAt:  time.Now(),
		Creater:    creater,
		UpdatedAt:  time.Now(),
	}
}
