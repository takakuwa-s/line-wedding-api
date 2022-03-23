package gateway

import (
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type FileRepository struct {
}

// Newコンストラクタ
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (fr *FileRepository) SaveFile(file *entity.File) error {
	conf.Log.Info("Successfully save the file", zap.Any("file", file))
	return nil
}