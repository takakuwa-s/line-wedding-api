package controller

import (
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"os"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
)

type LineController struct {
	bot *linebot.Client
	ml  *usecase.MessageHandler
}

// コンストラクタ
func NewLineController(ml *usecase.MessageHandler) *LineController {
	accessToken := os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		logger.Error("Failed to create lineBot instance", zap.Any("err", err))
	}
	return &LineController{bot: bot, ml: ml}
}

func (lw *LineController) Webhook(c *gin.Context) {
	events, err := lw.bot.ParseRequest(c.Request)
	if err != nil {
		logger.Error("Failed to parse the request", zap.Any("err", err))
		return
	}
	logger.Info("Successfully parse the request", zap.Any("events", events))
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			message := dto.Message{
				ReplyToken : event.ReplyToken,
				Text : event.Message.(*linebot.TextMessage).Text,
			}
			lw.ml.HandleTextMessage(message)
		}
	}
}