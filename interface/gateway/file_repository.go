package gateway

import (
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"go.uber.org/zap"
	"google.golang.org/api/iterator"

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
		return fmt.Errorf("failed adding a new file; file =  %v, err = %w", file, err)
	}
	conf.Log.Info("Successfully save the file", zap.Any("file", file))
	return nil
}

func (fr *FileRepository) DeleteFile(id, updater string) error {
	if _, err := fr.f.Client.Collection("files").Doc(id).Update(conf.Ctx, []firestore.Update{
		{
			Path:  "IsDeleted",
			Value: true,
		},
		{
			Path:  "UpdatedAt",
			Value: time.Now(),
		},
		{
			Path:  "Updater",
			Value: updater,
		},
	}); err != nil {
		return fmt.Errorf("failed delete the file; id =  %s, err = %w", id, err)
	}
	conf.Log.Info("Successfully delete the file", zap.String("id", id))
	return nil
}

func (fr *FileRepository) FindByCreaterAndIsDeleted(creater string, isDeleted bool) ([]entity.File, error) {
	var files []entity.File
	iter := fr.f.Client.Collection("files").Where("Creater", "==", creater).Where("IsDeleted", "==", isDeleted).Documents(conf.Ctx)
	for dsnap, err := iter.Next(); err != iterator.Done; dsnap, err = iter.Next() {
		if err != nil {
			return nil, fmt.Errorf("failed get a file; err = %w", err)
		}
		var f entity.File
		dsnap.DataTo(&f)
		files = append(files, f)
	}
	conf.Log.Info("Successfully find the files of the creater", zap.String("creater", creater))
	return files, nil
}
