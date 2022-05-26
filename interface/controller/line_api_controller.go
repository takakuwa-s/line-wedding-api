package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type LineApiController struct {
	au *usecase.ApiUsecase
	lpu *usecase.LinePushUsecase
}

// コンストラクタ
func NewLineApiController(
	au *usecase.ApiUsecase,
	lpu *usecase.LinePushUsecase) *LineApiController {
	return &LineApiController{au: au, lpu: lpu}
}

func (lac *LineApiController) SendMessageToLineBot(c *gin.Context) {
	err := lac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[SendMessageToLineBot] Authorization failed", zap.String("error", err.Error()))
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// return
	}
	messageKey := c.Query("messageKey")
	if err := lac.lpu.PublishMessageToAttendee(messageKey); err != nil {
		conf.Log.Error("[SendMessageToLineBot] Sending messages failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}