package dto

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type InitApiResponse struct {
	User  *entity.User   `json:"user"`
	Files []FileResponce `json:"files"`
}

func NewInitApiResponse(user *entity.User, files []FileResponce) *InitApiResponse {
	return &InitApiResponse{
		User:  user,
		Files: files,
	}
}
