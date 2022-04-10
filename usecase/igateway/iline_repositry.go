package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type ILineRepository interface {
	GetUserProfileById(id string, botType dto.BotType) (*entity.User, error)
	GetQuotaComsuption(botType dto.BotType) (int64, error)
	GetFileContent(botType dto.BotType, messageId string) (io.ReadCloser, error)
}
