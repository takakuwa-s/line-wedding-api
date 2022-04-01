package gateway

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type FileRepository struct {
	f *dto.Firestore
}

// Newコンストラクタ
func NewFileRepository(f *dto.Firestore) *FileRepository {
	return &FileRepository{f:f}
}

func (fr *FileRepository) SaveFile(file *entity.File) error {
	if _, err := fr.f.Client.Collection("files").Doc(file.Id).Set(fr.f.Ctx, file); err != nil {
		return fmt.Errorf("failed adding a new file; file =  %v, err = %w", file, err)
	}
	conf.Log.Info("Successfully save the file", zap.Any("file", file))
	return nil
}