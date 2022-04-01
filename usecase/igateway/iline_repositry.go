package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type ILineRepository interface {
	GetUserProfileById(id string, botType dto.BotType) (*entity.User, error)
	GetQuotaComsuption(botType dto.BotType) (int64, error)
}
