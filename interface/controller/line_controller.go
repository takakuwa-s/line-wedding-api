package controller

import (
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"log"
	"net/http"
	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
  "github.com/gin-gonic/gin"
)

type LineController struct {
	bot *linebot.Client
	ml *usecase.MessageHandler
}

// コンストラクタ
func NewLineController(ml *usecase.MessageHandler) *LineController {
	accessToken := os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
	}
	return &LineController{bot:bot, ml:ml}
}

func (lw *LineController) Webhook(c *gin.Context) {
	events, err := lw.bot.ParseRequest(c.Request)
	if err != nil {
    c.String(http.StatusInternalServerError, "System error when parsing the request; err = %v", err)
		return
	}
  log.Printf("Successfully parse the request; events = %v", events)
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