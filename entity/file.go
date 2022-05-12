package entity

import (
	"time"
)

type File struct {
	LineFileId string
	FileId string
	Name string
	ContentUri string
	ThumbnailUri string
	IsUploaded bool
	CreatedAt time.Time
	Creater string
	UpdatedAt time.Time
	Updater string
}

func NewFile(messageId, creater string) *File {
	return &File{
		LineFileId: messageId,
		IsUploaded: false,
		CreatedAt: time.Now(),
		Creater: creater,
		UpdatedAt: time.Now(),
	}
}