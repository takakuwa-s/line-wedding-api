package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IUserRepository interface {
	SaveUser(user *entity.User) error
	UpdateFollowStatusById(id string, status bool) error
	FindByWillJoinAndFollowStatus(willJoin, followStatus bool) (*[]entity.User, error)
	FindByIsAdmin(isAdmin bool) (*[]entity.User, error)
	FindById(id string) (*entity.User, error)
}
