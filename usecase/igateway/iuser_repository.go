package igateway

import "github.com/takakuwa-s/line-wedding-api/usecase/dto"

type IUserRepository interface {
	SaveUser(user *dto.User)
	UpdateFollowStatusById(id string, status bool)
}