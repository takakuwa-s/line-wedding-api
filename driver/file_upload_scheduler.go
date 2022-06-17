package driver

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

type FileUploadScheduler struct {
	fuu *usecase.FileUploadUsecase
}

// Newコンストラクタ
func NewFileUploadScheduler(fuu *usecase.FileUploadUsecase) *FileUploadScheduler {
	return &FileUploadScheduler{fuu: fuu}
}

// Init ルーティング設定
func (fus *FileUploadScheduler) Init() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Minute().Do(fus.fuu.RecoverFileUploading)
	scheduler.StartAsync()
}
