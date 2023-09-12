package dto

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type UpdateUserRequest struct {
	Id         string
	Attendance bool
	GuestType  string
	Name       string
	NameKana   string
	Email      string
	PostalCode string
	Address    string
	Allergy    string
	Companions []entity.Companion
	Message    string
}

func (ur *UpdateUserRequest) ToUser(u *entity.User) *entity.User {
	u.Attendance = ur.Attendance
	u.GuestType = ur.GuestType
	u.Name = ur.Name
	u.NameKana = ur.NameKana
	u.Email = ur.Email
	u.PostalCode = ur.PostalCode
	u.Address = ur.Address
	u.Allergy = ur.Allergy
	u.Companions = ur.Companions
	u.Message = ur.Message
	return u
}
