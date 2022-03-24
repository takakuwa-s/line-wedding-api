package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type ILineRepository interface {
	FindUserById(id string, botType dto.BotType) (*entity.User, error)
	GetQuotaComsuption(botType dto.BotType) (int, error)
}
