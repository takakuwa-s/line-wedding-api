package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

type IBinaryRepository interface {
	SaveImageBinary(file entity.File, content io.ReadCloser) (*entity.File, error)
	SaveVideoBinary(file entity.File, content io.ReadCloser) (*entity.File, error)
	SaveSlideShowBinary(s entity.SlideShow, content io.Reader) (*entity.SlideShow, error)
	DeleteBinary(name, prefix string) error
}
