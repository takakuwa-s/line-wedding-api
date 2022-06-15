package gateway

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type Entity interface {
}

type CommonRepository struct {
	f *dto.Firestore
}

// Newコンストラクタ
func NewCommonRepository(f *dto.Firestore) *CommonRepository {
	return &CommonRepository{f: f}
}

func (cr *CommonRepository) Save(collectionName string, id string, data interface{}) error {
	if _, err := cr.f.Firestore.Collection(collectionName).Doc(id).Set(conf.Ctx, data); err != nil {
		return fmt.Errorf("failed to save data; collectionName = %s, data =  %v, err = %w", collectionName, data, err)
	}
	conf.Log.Info("Successfully save the data", zap.String("collectionName", collectionName), zap.Any("data", data))
	return nil
}

func (cr *CommonRepository) DeleteById(collectionName, id string) error {
	_, err := cr.f.Firestore.Collection(collectionName).Doc(id).Delete(conf.Ctx)
	if err != nil {
		return fmt.Errorf("failed to delete the data from firestore; collectionName = %s, id =  %s, err = %w", collectionName, id, err)
	}
	conf.Log.Info("Successfully delete the data from firestore", zap.String("collectionName", collectionName), zap.String("id", id))
	return nil
}

func (cr *CommonRepository) DeleteByIds(collectionName string, ids []string) error {
	for _, id := range ids {
		err := cr.DeleteById(collectionName, id)
		if err != nil {
			return err
		}
	}
	conf.Log.Info("Successfully delete the multiple data", zap.String("collectionName", collectionName), zap.Strings("ids", ids))
	return nil
}

func (cr *CommonRepository) FindById(collectionName, id string) (*firestore.DocumentSnapshot, error) {
	dsnap, err := cr.f.Firestore.Collection(collectionName).Doc(id).Get(conf.Ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to find the data from firestore; collectionName = %s, id =  %s, err = %w", collectionName, id, err)
		}
	}
	conf.Log.Info("Successfully find the data by Id", zap.String("collectionName", collectionName), zap.String("id", id), zap.Any("dsnap", dsnap))
	return dsnap, nil
}

func (cr *CommonRepository) SplitSlice(ids []string) [][]string {
	limit := 4
	if len(ids) <= limit {
		return [][]string{ids}
	}
	var result [][]string
	for i := 0; i < len(ids); i += limit {
		if i+limit < len(ids) {
			result = append(result, ids[i:i+limit])
		} else {
			result = append(result, ids[i:])
		}
	}
	return result
}
