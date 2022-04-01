package entity

import (
	"time"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type User struct {
	Id string
	Name string
	WillJoin bool
	IsAdmin bool
	FollowStatus bool
	Language string
	IconUrl string
	StatusMessage string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(res *linebot.UserProfileResponse) *User {
	return &User{
		Id: res.UserID,
		Name: res.DisplayName,
		WillJoin: false,
		IsAdmin: false,
		FollowStatus: true,
		Language: res.Language,
		IconUrl: res.PictureURL,
		StatusMessage: res.StatusMessage,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}