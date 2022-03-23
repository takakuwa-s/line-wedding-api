package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IUserRepository interface {
	SaveUser(user *entity.User) error
	UpdateFollowStatusById(id string, status bool) error
	FindByWillJoin(willJoin bool) (*[]entity.User, error)
}
