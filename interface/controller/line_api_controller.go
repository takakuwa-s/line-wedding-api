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
}

// コンストラクタ
func NewLineApiController(
	au *usecase.ApiUsecase) *LineApiController {
	return &LineApiController{au: au}
}

func (lac *LineApiController) SendMessageToLineBot(c *gin.Context) {
	messageKey := c.Query("messageKey")
	flag := c.Query("flag")
	if messageKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "messageKey is required"})
		return
	}
	if flag != "Attendance" && flag != "Follow" && flag != "IsAdmin" && flag != "Registered" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "flag must be one of Attendance, Follow, IsAdmin, Registered"})
		return
	}
	if err := lac.au.PublishMessageToUsers(messageKey, flag); err != nil {
		conf.Log.Error("[SendMessageToLineBot] Sending messages failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
