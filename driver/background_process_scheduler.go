package driver

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

type BackgroundProcessScheduler struct {
	bpu *usecase.BackgroundProcessUsecase
}

// Newコンストラクタ
func NewBackgroundProcessScheduler(bpu *usecase.BackgroundProcessUsecase) *BackgroundProcessScheduler {
	return &BackgroundProcessScheduler{bpu: bpu}
}

// Init ルーティング設定
func (bps *BackgroundProcessScheduler) Init() {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(5).Minute().Do(bps.bpu.RecoverFileProcessing)
	scheduler.StartAsync()
}
