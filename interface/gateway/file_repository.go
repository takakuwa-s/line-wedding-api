package gateway

import (
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type FileRepository struct {
}

// Newコンストラクタ
func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (fr *FileRepository) SaveFile(file *dto.File) {
	conf.Log.Info("Successfully save the file", zap.Any("file", file))
}