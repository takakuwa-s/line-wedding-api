package igateway

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type ISlideShowRepository interface {
	SaveSlideShow(s *entity.SlideShow) error
	UpdateSelectedById(selected bool, id string) error
	DeleteById(id string) error
	FindById(id string) (*entity.SlideShow, error)
	FindAllOrderByUpdatedAt() ([]entity.SlideShow, error)
	FindBySelectedOrderByUpdatedAt(selected bool) ([]entity.SlideShow, error)
}
