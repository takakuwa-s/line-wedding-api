package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IUserRepository interface {
	SaveUser(user *entity.User) error
	UpdateFollowById(id string, follow bool) error
	FindById(id string) (*entity.User, error)
	FindByIsAdmin(isAdmin bool) ([]entity.User, error)
	FindByAttendanceAndFollow(attendance, follow bool) ([]entity.User, error)
	FindByFlagOrderByName(limit int, startId, flag string, val bool) ([]entity.User, error)
}
