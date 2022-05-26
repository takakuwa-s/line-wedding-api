package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

type ILineGateway interface {
	GetUserProfileById(id string) (*entity.User, error)
	GetQuotaComsuption() (int64, error)
	GetFileContent(messageId string) (io.ReadCloser, error)
}
