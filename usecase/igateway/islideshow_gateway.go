package igateway

import (
	"io"

	"github.com/takakuwa-s/line-wedding-api/dto"
)

type ISlideShowGateway interface {
	Render(imageUrls, videoUrls []string) (*dto.SlideShowCreateResponce, error)
	DownloadContent(url string) (io.Reader, error)
}
