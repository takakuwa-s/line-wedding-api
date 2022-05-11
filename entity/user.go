package entity

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type User struct {
	Id string
	LineName string
	FollowStatus bool
	Language string
	IconUrl string
	StatusMessage string
	Attendance bool
  GuestType string
  FamilyName string
  FirstName string
  FamilyNameKana string
  FirstNameKana string
  PhoneNumber string
  PostalCode string
  Address string
  Allergy string
  Message string
	IsAdmin bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(res *linebot.UserProfileResponse) *User {
	return &User{
		Id: res.UserID,
		LineName: res.DisplayName,
		IsAdmin: false,
		FollowStatus: true,
		Language: res.Language,
		IconUrl: res.PictureURL,
		StatusMessage: res.StatusMessage,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}