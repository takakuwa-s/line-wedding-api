package entity

import (
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type Companion struct {
	Name    string `json:"name"`
	Allergy string `json:"allergy"`
}

type User struct {
	Id         string      `json:"id"`
	LineName   string      `json:"lineName"`
	Follow     bool        `json:"follow"`
	Language   string      `json:"language"`
	IconUrl    string      `json:"iconUrl"`
	Attendance bool        `json:"attendance"`
	GuestType  string      `json:"guestType"`
	Name       string      `json:"name"`
	NameKana   string      `json:"nameKana"`
	Email      string      `json:"email"`
	PostalCode string      `json:"postalCode"`
	Address    string      `json:"address"`
	Allergy    string      `json:"allergy"`
	Message    string      `json:"message"`
	Companions []Companion `json:"companions"`
	IsAdmin    bool        `json:"isAdmin"`
	Registered bool        `json:"registered"`
	Note       string      `json:"note"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
}

func NewUser(res *linebot.UserProfileResponse) *User {
	return &User{
		Id:         res.UserID,
		LineName:   res.DisplayName,
		Attendance: true,
		IsAdmin:    false,
		Registered: false,
		Follow:     true,
		Language:   res.Language,
		IconUrl:    res.PictureURL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
