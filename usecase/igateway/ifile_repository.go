package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type IFileRepository interface {
	SaveFile(file *dto.File)
}
