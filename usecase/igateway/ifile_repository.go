package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IFileRepository interface {
	SaveFile(file *entity.File) error
	DeleteFileById(id string) error
	DeleteFileByIds(ids []string) error
	FindById(id string) (*entity.File, error)
	FindByIds(ids []string) ([]entity.File, error)
	FindByIdsAndUploaded(ids []string, uploaded bool) ([]entity.File, error)
	FindByLimitAndStartIdAndUserIdAndUploaded(limit int, startId, userId, orderBy string, uploaded *bool) ([]entity.File, error)
	FindByUploadedOrCalculatedFalse() ([]entity.File, error)
}
