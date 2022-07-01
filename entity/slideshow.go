package entity

import (
	"time"
)

type SlideShow struct {
	Id           string    `json:"id"`
	Name         string    `json:"name"`
	ContentUrl   string    `json:"contentUrl"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
	MimeType     string    `json:"mimeType"`
	Selected     bool      `json:"selected"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewSlideShow(id string) *SlideShow {
	return &SlideShow{
		Id:        id,
		Selected:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
