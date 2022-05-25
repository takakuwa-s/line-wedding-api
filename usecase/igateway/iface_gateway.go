package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type IFaceGateway interface {
	GetFaceAnalysis(url string) ([]*dto.FaceResponse, error)
}
