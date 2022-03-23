package controller

import (
	"github.com/takakuwa-s/line-wedding-api/usecase"
)

type ScheduleController struct {
	pmu *usecase.PushMessageUsecase
}

// コンストラクタ
func NewScheduleController(pmu *usecase.PushMessageUsecase) *ScheduleController {
	return &ScheduleController{pmu: pmu}
}