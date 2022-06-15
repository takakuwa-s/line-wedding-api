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
	cr *CommonRepository
	f  *dto.Firestore
}

// Newコンストラクタ
func NewFileRepository(cr *CommonRepository, f *dto.Firestore) *FileRepository {
	return &FileRepository{cr: cr, f: f}
}

func (fr *FileRepository) SaveFile(file *entity.File) error {
	return fr.cr.Save("files", file.Id, file)
}

func (fr *FileRepository) DeleteFileById(id string) error {
	return fr.cr.DeleteById("files", id)
}

func (fr *FileRepository) DeleteFileByIds(ids []string) error {
	return fr.cr.DeleteByIds("files", ids)
}

func (fr *FileRepository) FindById(id string) (*entity.File, error) {
	dsnap, err := fr.cr.FindById("files", id)
	if err != nil {
		return nil, err
	}
	if dsnap == nil {
		return nil, nil
	}
	var file entity.File
	dsnap.DataTo(&file)
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
	var files []entity.File
	for _, list := range fr.cr.SplitSlice(ids) {
		query := fr.f.Firestore.Collection("files").Where("Id", "in", list)
		f, err := fr.executeQuery(&query)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Strings("ids", ids))
	return files, nil
}

func (fr *FileRepository) FindByIdsAndUploaded(ids []string, uploaded bool) ([]entity.File, error) {
	var files []entity.File
	for _, list := range fr.cr.SplitSlice(ids) {
		query := fr.f.Firestore.Collection("files").Where("Id", "in", list).Where("Uploaded", "==", uploaded)
		f, err := fr.executeQuery(&query)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Strings("ids", ids), zap.Bool("uploaded", uploaded))
	return files, nil
}

func (fr *FileRepository) FindByLimitAndStartIdAndUserIdAndUploaded(limit int, startId, userId, orderBy string, uploaded *bool) ([]entity.File, error) {
	conf.Log.Info("3", zap.Any("uploaded", &uploaded))
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
	if uploaded != nil {
		query = query.Where("Uploaded", "==", &uploaded)
	}
	files, err := fr.executeQuery(&query)
	if err != nil {
		return nil, err
	}
	conf.Log.Info("Successfully find the file metadata with", zap.Int("file count", len(files)), zap.Int("limit", limit), zap.String("startId", startId), zap.String("userId", userId))
	return files, nil
}
