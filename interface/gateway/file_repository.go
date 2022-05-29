package gateway

import (
	"fmt"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type FileRepository struct {
	f *dto.Firestore
}

// Newコンストラクタ
func NewFileRepository(f *dto.Firestore) *FileRepository {
	return &FileRepository{f: f}
}

func (fr *FileRepository) SaveFile(file *entity.File) error {
	if _, err := fr.f.Firestore.Collection("files").Doc(file.Id).Set(conf.Ctx, file); err != nil {
		return fmt.Errorf("failed to save a new file metadata; file =  %v, err = %w", file, err)
	}
	conf.Log.Info("Successfully save the file metadata", zap.Any("file", file))
	return nil
}

func (fr *FileRepository) DeleteFileById(id string) error {
	_, err := fr.f.Firestore.Collection("files").Doc(id).Delete(conf.Ctx)
	if err != nil {
		return fmt.Errorf("failed to delete the file metadata; id =  %s, err = %w", id, err)
	}
	conf.Log.Info("Successfully delete the file metadata", zap.String("id", id))
	return nil
}

func (fr *FileRepository) FindById(id string) (*entity.File, error) {
	dsnap, err := fr.f.Firestore.Collection("files").Doc(id).Get(conf.Ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return nil, fmt.Errorf("failed to find the file metadata; id =  %s, err = %w", id, err)
		}
	}
	var file entity.File
	dsnap.DataTo(&file)
	conf.Log.Info("Successfully find the file metadata by Id", zap.String("id", id), zap.Any("file", file))
	return &file, nil
}

func (fr *FileRepository) executeQuery(query *firestore.Query) ([]entity.File, error) {
	var files []entity.File
	iter := query.Documents(conf.Ctx)
	for {
		dsnap, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get a file metadata ; err = %w", err)
		}
		var f entity.File
		dsnap.DataTo(&f)
		files = append(files, f)
	}
	if files == nil {
		return []entity.File{}, nil
	}
	return files, nil
}

func (fr *FileRepository) FindByIds(ids []string) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").Where("Id", "in", ids)
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Strings("ids", ids))
	return files, nil
}

func (fr *FileRepository) FindByLimitAndStartIdAndUserId(limit int, startId, userId, orderBy string) ([]entity.File, error) {
	query := fr.f.Firestore.Collection("files").OrderBy(orderBy, firestore.Desc)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if startId != "" {
		dsnap, err := fr.f.Firestore.Collection("files").Doc(startId).Get(conf.Ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return []entity.File{}, nil
			} else {
				return nil, fmt.Errorf("failed to get the file metadata by startId; id =  %s err = %w", startId, err)
			}
		}
		query = query.StartAfter(dsnap)
	}
	if userId != "" {
		query = query.Where("Creater", "==", userId)
	}
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Int("limit", limit), zap.String("startId", startId), zap.String("userId", userId))
	return files, nil
}
