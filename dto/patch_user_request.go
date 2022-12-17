package dto

import (
	"fmt"
	"strconv"
)

type PatchUserRequest struct {
	IsAdmin    string `json:"isAdmin"`
	Attendance string `json:"attendance"`
	Note       string `json:"note"`
}

func (p *PatchUserRequest) GetFieldAndVal() (string, interface{}, error) {
	if p.IsAdmin != "" && p.Attendance == "" && p.Note == "" {
		if val, err := strconv.ParseBool(p.IsAdmin); err == nil {
			return "IsAdmin", val, nil
		} else {
			return "", false, fmt.Errorf("[PatchUser] isAdmin must be boolean")
		}
	} else if p.IsAdmin == "" && p.Attendance != "" && p.Note == "" {
		if val, err := strconv.ParseBool(p.Attendance); err == nil {
			return "Attendance", val, nil
		} else {
			return "", false, fmt.Errorf("[PatchUser] attendance must be boolean")
		}
	} else if p.IsAdmin == "" && p.Attendance == "" && p.Note != "" {
		return "Note", p.Note, nil
	} else {
		return "", false, fmt.Errorf("[PatchUser] One field can be updated at a time")
	}
}
