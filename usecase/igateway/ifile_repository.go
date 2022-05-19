package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IFileRepository interface {
	SaveFile(file *entity.File) error
	DeleteFile(id string) error
	FindById(id string) (*entity.File, error)
	FindByLimitAndStartIdAndUserId(limit int, startId, userId string) ([]entity.File, error)
}
