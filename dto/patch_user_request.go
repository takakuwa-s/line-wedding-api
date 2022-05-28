package dto

import (
	"fmt"
	"strconv"
)

type PatchUserRequest struct {
	IsAdmin    string `json:"isAdmin"`
	Attendance string `json:"attendance"`
}

func (p *PatchUserRequest) GetFieldAndVal() (string, bool, error) {
	if p.IsAdmin != "" && p.Attendance != "" {
		return "", false, fmt.Errorf("[PatchUser] isAdmin and attendance are exclusive")
	} else if p.IsAdmin == "" && p.Attendance == "" {
		return "", false, fmt.Errorf("[PatchUser] neither isAdmin nor attendance are set")
	} else if p.IsAdmin != "" {
		if val, err := strconv.ParseBool(p.IsAdmin); err == nil {
			return "IsAdmin", val, nil
		} else {
			return "", false, fmt.Errorf("[PatchUser] isAdmin must be boolean")
		}
	} else {
		if val, err := strconv.ParseBool(p.Attendance); err == nil {
			return "Attendance", val, nil
		} else {
			return "", false, fmt.Errorf("[PatchUser] attendance must be boolean")
		}
	}
}
