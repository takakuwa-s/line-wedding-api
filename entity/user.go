package entity

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type User struct {
	Id string `json:"id"`
	LineName string `json:"lineName"`
	FollowStatus bool `json:"followStatus"`
	Language string `json:"language"`
	IconUrl string `json:"iconUrl"`
	StatusMessage string `json:"statusMessage"`
	Attendance bool `json:"attendance"`
  GuestType string `json:"guestType"`
  FamilyName string `json:"familyName"`
  FirstName string `json:"firstName"`
  FamilyNameKana string `json:"familyNameKana"`
  FirstNameKana string `json:"firstNameKana"`
  PhoneNumber string `json:"phoneNumber"`
  PostalCode string `json:"postalCode"`
  Address string `json:"address"`
  Allergy string `json:"allergy"`
  Message string `json:"message"`
	IsAdmin bool `json:"isAdmin"`
	IsRegistered bool `json:"isRegistered"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(res *linebot.UserProfileResponse) *User {
	return &User{
		Id: res.UserID,
		LineName: res.DisplayName,
		IsAdmin: false,
		IsRegistered: false,
		FollowStatus: true,
		Language: res.Language,
		IconUrl: res.PictureURL,
		StatusMessage: res.StatusMessage,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}