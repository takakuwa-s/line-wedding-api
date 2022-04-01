package entity

import (
	"time"
)

type File struct {
	Id string
	IsUploaded bool
	CreatedAt time.Time
	CreatedBy string
	UpdatedAt time.Time
}

func NewFile(messageId, createdBy string) *File {
	return &File{
		Id: messageId,
		IsUploaded: false,
		CreatedAt: time.Now(),
		CreatedBy: createdBy,
		UpdatedAt: time.Now(),
	}
}