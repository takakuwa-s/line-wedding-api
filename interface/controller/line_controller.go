package controller

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type LineController struct {
	bot *linebot.Client
	ml  *usecase.MessageHandler
}

// コンストラクタ
func NewLineController(bot *linebot.Client, ml *usecase.MessageHandler) *LineController {
	return &LineController{bot: bot, ml: ml}
}

func (lw *LineController) Webhook(c *gin.Context) {
	events, err := lw.bot.ParseRequest(c.Request)
	if err != nil {
		conf.Log.Error("Failed to parse the request", zap.Any("err", err))
		return
	}
	conf.Log.Info("Successfully parse the request", zap.Any("events", events))
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			message := dto.RequestMessage{
				ReplyToken : event.ReplyToken,
				Text : event.Message.(*linebot.TextMessage).Text,
			}
			lw.ml.HandleTextMessage(message)
		}
	}
}