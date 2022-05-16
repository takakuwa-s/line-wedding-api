package entity

import (
	"time"
)

type FileType string

const (
    ImageType FileType = FileType("image")
    VideoType FileType = FileType("video")
)

func (f *FileType) ToString() string {
	return string(*f)
}

type File struct {
	LineFileId string `json:"lineFileId"`
	FileId string `json:"fileId"`
	Name string `json:"name"`
	FileType FileType `json:"fileType"`
	ContentUrl string `json:"contentUrl"`
	ThumbnailUrl string `json:"thumbnailUrl"`
	Width int64 `json:"width"`
	Height int64 `json:"height"`
	MimeType string `json:"mimeType"`
	IsUploaded bool `json:"isUploaded"`
	CreatedAt time.Time `json:"createdAt"`
	Creater string `json:"creater"`
	UpdatedAt time.Time `json:"updatedAt"`
	Updater string `json:"updater"`
}

func NewFile(messageId, creater string, fileType FileType) *File {
	return &File{
		LineFileId: messageId,
		FileType: fileType,
		IsUploaded: false,
		CreatedAt: time.Now(),
		Creater: creater,
		UpdatedAt: time.Now(),
	}
}