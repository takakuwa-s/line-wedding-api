package entity

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type User struct {
	Id string `json:"id"`
	LineName string `json:"lineName"`
	Follow bool `json:"follow"`
	Language string `json:"language"`
	IconUrl string `json:"iconUrl"`
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
	Registered bool `json:"registered"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(res *linebot.UserProfileResponse) *User {
	return &User{
		Id: res.UserID,
		LineName: res.DisplayName,
		IsAdmin: false,
		Registered: false,
		Follow: true,
		Language: res.Language,
		IconUrl: res.PictureURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}