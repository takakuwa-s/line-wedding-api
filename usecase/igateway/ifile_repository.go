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
	FindByLimitAndStartIdAndUserIdAndFileTypeAndUploaded(limit int, startId, userId, orderBy, fileType string, uploaded *bool) ([]entity.File, error)
	FindByUploadedOrCalculatedFalse() ([]entity.File, error)
	FindByUploadedAndFileType(limit int, uploaded bool, fileType entity.FileType) ([]entity.File, error)
	FindByUploadedAndFileTypeAndDuration(limit int, uploaded bool, fileType entity.FileType, duration int) ([]entity.File, error)
}
