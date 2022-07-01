package gateway

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
)

type SlideShowRepository struct {
	cr *CommonRepository
	f  *dto.Firestore
}

// Newコンストラクタ
func NewFSlideShowRepository(cr *CommonRepository, f *dto.Firestore) *SlideShowRepository {
	return &SlideShowRepository{cr: cr, f: f}
}

func (sr *SlideShowRepository) SaveSlideShow(s *entity.SlideShow) error {
	return sr.cr.Save("slideshows", s.Id, s)
}

func (sr *SlideShowRepository) DeleteById(id string) error {
	return sr.cr.DeleteById("slideshows", id)
}

func (sr *SlideShowRepository) FindById(id string) (*entity.SlideShow, error) {
	dsnap, err := sr.cr.FindById("slideshows", id)
	if err != nil {
		return nil, err
	}
	if dsnap == nil {
		return nil, nil
	}
	var s entity.SlideShow
	dsnap.DataTo(&s)
	return &s, nil
}

func (sr *SlideShowRepository) UpdateSelectedById(selected bool, id string) error {
	if _, err := sr.f.Firestore.Collection("slideshows").Doc(id).Update(conf.Ctx, []firestore.Update{
		{
			Path:  "Selected",
			Value: selected,
		},
		{
			Path:  "UpdatedAt",
			Value: time.Now(),
		},
	}); err != nil {
		return fmt.Errorf("failed to update the slideshow; id =  %s, selected = %t, err = %w", id, selected, err)
	}
	conf.Log.Info("Successfully update the slideshow", zap.String("id", id), zap.Bool("selected", selected))
	return nil
}

func (sr *SlideShowRepository) executeQuery(query *firestore.Query) ([]entity.SlideShow, error) {
	var list []entity.SlideShow
	iter := query.Documents(conf.Ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get a slide show metadata ; err = %w", err)
		}
		var s entity.SlideShow
		dsnap.DataTo(&s)
		list = append(list, s)
	}
	if list == nil {
		return []entity.SlideShow{}, nil
	}
	return list, nil
}

func (sr *SlideShowRepository) FindAllOrderByUpdatedAt() ([]entity.SlideShow, error) {
	query := sr.f.Firestore.Collection("slideshows").OrderBy("UpdatedAt", firestore.Desc)
	res, err := sr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (sr *SlideShowRepository) FindBySelectedOrderByUpdatedAt(selected bool) ([]entity.SlideShow, error) {
	query := sr.f.Firestore.Collection("slideshows").Where("Selected", "==", selected).OrderBy("UpdatedAt", firestore.Desc)
	res, err := sr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
