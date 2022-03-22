package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type ILineRepository interface {
	FindUserById(id string) *dto.User
}
