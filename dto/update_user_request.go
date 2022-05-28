package dto

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type UpdateUserRequest struct {
	Id             string
	Attendance     bool
	GuestType      string
	FamilyName     string
	FirstName      string
	FamilyNameKana string
	FirstNameKana  string
	PhoneNumber    string
	PostalCode     string
	Address        string
	Allergy        string
	Message        string
}

func (ur *UpdateUserRequest) ToUser(u *entity.User) *entity.User {
	u.Attendance = ur.Attendance
	u.GuestType = ur.GuestType
	u.FamilyName = ur.FamilyName
	u.FirstName = ur.FirstName
	u.FamilyNameKana = ur.FamilyNameKana
	u.FirstNameKana = ur.FirstNameKana
	u.PhoneNumber = ur.PhoneNumber
	u.PostalCode = ur.PostalCode
	u.Address = ur.Address
	u.Allergy = ur.Allergy
	u.Message = ur.Message
	return u
}
