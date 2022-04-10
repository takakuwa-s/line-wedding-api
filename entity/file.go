package entity

import (
	"time"
)

type File struct {
	LineFileId string
	FileId string
	Name string
	Uri string
	IsUploaded bool
	CreatedAt time.Time
	Creater string
	UpdatedAt time.Time
	Updater string
	IsDeleted bool
}

func NewFile(messageId, creater string) *File {
	return &File{
		LineFileId: messageId,
		IsUploaded: false,
		CreatedAt: time.Now(),
		Creater: creater,
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}
}