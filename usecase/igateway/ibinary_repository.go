package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IBinaryRepository interface {
	SaveImageBinary(file entity.File, content io.ReadCloser) (*entity.File, error)
	SaveVideoBinary(file entity.File, content io.ReadCloser) (*entity.File, error)
	DeleteBinary(name string, fileType entity.FileType) error
}
