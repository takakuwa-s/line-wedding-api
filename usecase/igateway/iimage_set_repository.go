package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IImageSetRepository interface {
	DeleteById(id string) error
	AppendFileIdByImageSet(set *entity.ImageSet, fileId string) (*entity.ImageSet, error)
}
