package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IFileRepository interface {
	SaveFile(file *entity.File) error
	DeleteFile(id, updatedBy string) error
	FindByCreaterAndIsDeleted(creater string, isDeleted bool) ([]entity.File, error)
}
