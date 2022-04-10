package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IBinaryRepository interface {
	SaveBinary(file *entity.File, content io.ReadCloser) (*entity.File, error)
	DeleteBinary(id string) error
}
