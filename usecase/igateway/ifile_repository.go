package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IFileRepository interface {
	SaveFile(file *entity.File) error
	DeleteFile(id string) error
	FindById(id string) (*entity.File, error)
	FindByLimit(limit int) ([]entity.File, error)
	FindByLimitAndStartId(limit int, startId string) ([]entity.File, error)
}
