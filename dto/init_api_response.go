package dto

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type InitApiResponse struct {
	User  *entity.User  `json:"user"`
	Files []entity.File `json:"files"`
}

func NewInitApiResponse(user *entity.User, files []entity.File) *InitApiResponse {
	return &InitApiResponse{
		User:  user,
		Files: files,
	}
}
