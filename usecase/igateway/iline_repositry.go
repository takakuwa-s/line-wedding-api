package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type ILineRepository interface {
	FindUserById(id string) (*entity.User, error)
	GetQuotaComsuption() (int, error)
}
