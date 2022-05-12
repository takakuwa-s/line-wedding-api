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
	if _, err := fr.f.Client.Collection("files").Doc(file.LineFileId).Set(conf.Ctx, file); err != nil {
		return fmt.Errorf("failed to save a new file metadata; file =  %v, err = %w", file, err)
	}
	conf.Log.Info("Successfully save the file metadata", zap.Any("file", file))
	return nil
}

func (fr *FileRepository) DeleteFile(id string) error {
	_, err := fr.f.Client.Collection("files").Doc(id).Delete(conf.Ctx)
	if err != nil {
		return fmt.Errorf("failed to delete the file metadata; id =  %s, err = %w", id, err)
	}
	conf.Log.Info("Successfully delete the file metadata", zap.String("id", id))
	return nil
}

func (fr *FileRepository)	FindById(id string) (*entity.File, error) {
	dsnap, err := fr.f.Client.Collection("files").Doc(id).Get(conf.Ctx)
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
	for dsnap, err := iter.Next(); err != iterator.Done; dsnap, err = iter.Next() {
		conf.Log.Info("dsnap", zap.Any("dsnap", dsnap))
		if err != nil {
			return nil, fmt.Errorf("failed to get a file metadata ; err = %w", err)
		}
		var f entity.File
		dsnap.DataTo(&f)
		files = append(files, f)
	}
	return files, nil
}

func (fr *FileRepository) FindByLimit(limit int) ([]entity.File, error) {
	query := fr.f.Client.Collection("files").
							OrderBy("CreatedAt", firestore.Desc).
							Limit(limit)
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("limit", limit))
	return files, nil
}

func (fr *FileRepository) FindByLimitAndStartId(limit int, startId string) ([]entity.File, error) {
	dsnap, err := fr.f.Client.Collection("files").Doc(startId).Get(conf.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the file metadata; err = %w", err)
	}
	query := fr.f.Client.Collection("files").
							OrderBy("CreatedAt", firestore.Desc).
							StartAfter(dsnap).
							Limit(limit)
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("limit", limit), zap.String("startId", startId))
	return files, nil
}
