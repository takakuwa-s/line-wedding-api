package entity

import (
	"time"
)

type File struct {
	Id string
	IsUploaded bool
	CreatedAt time.Time
	CreatedBy string
}

func NewFile(messageId, createdBy string, createdAt time.Time) *File {
	return &File{
		Id: messageId,
		IsUploaded: false,
		CreatedAt: createdAt,
		CreatedBy: createdBy,
	}
}