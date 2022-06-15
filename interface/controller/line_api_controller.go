package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type LineApiController struct {
	lpu *usecase.LinePushUsecase
}

// コンストラクタ
func NewLineApiController(
	lpu *usecase.LinePushUsecase) *LineApiController {
	return &LineApiController{lpu: lpu}
}

func (lac *LineApiController) SendMessageToLineBot(c *gin.Context) {
	messageKey := c.Query("messageKey")
	if messageKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messageKey is required"})
		return
	}
	if err := lac.lpu.PublishMessageToAttendee(messageKey); err != nil {
		conf.Log.Error("[SendMessageToLineBot] Sending messages failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
