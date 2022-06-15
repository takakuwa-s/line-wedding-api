package gateway

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ImageSetRepository struct {
	cr *CommonRepository
	f  *dto.Firestore
}

// Newコンストラクタ
func NewImageSetRepository(cr *CommonRepository, f *dto.Firestore) *ImageSetRepository {
	return &ImageSetRepository{cr: cr, f: f}
}

func (isr *ImageSetRepository) DeleteById(id string) error {
	return isr.cr.DeleteById("imageSets", id)
}

func (isr *ImageSetRepository) AppendFileIdByImageSet(set *entity.ImageSet, fileId string) (*entity.ImageSet, error) {
	ref := isr.f.Firestore.Collection("imageSets").Doc(set.Id)
	err := isr.f.Firestore.RunTransaction(conf.Ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				set.FileIds = append(set.FileIds, fileId)
				return tx.Create(ref, set)
			} else {
				return fmt.Errorf("failed to find the imageSets data in the transaction; id =  %s, err = %w", set.Id, err)
			}
		}
		fileIds, err := doc.DataAt("FileIds")
		if err != nil {
			return fmt.Errorf("failed to get the FileIds in the transaction; id =  %s, err = %w", set.Id, err)
		}
		fileIds = append(fileIds.([]interface{}), fileId)
		return tx.Set(ref, map[string]interface{}{
			"FileIds": fileIds,
		}, firestore.MergeAll)
	})
	if err != nil {
		return nil, err
	}
	dsnap, err := ref.Get(conf.Ctx)
	if err != nil {
		return nil, err
	}
	var res entity.ImageSet
	dsnap.DataTo(&res)
	return &res, err
}
