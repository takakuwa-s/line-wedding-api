package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IUserRepository interface {
	SaveUser(user *entity.User) error
	UpdateFollowStatusById(id string, status bool) error
	FindById(id string) (*entity.User, error)
	FindByIsAdmin(isAdmin bool) ([]entity.User, error)
	FindByAttendanceAndFollowStatus(attendance, followStatus bool) ([]entity.User, error)
}
