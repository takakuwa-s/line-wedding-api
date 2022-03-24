package controller

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/admin"
	"github.com/takakuwa-s/line-wedding-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type AdminLineController struct {
	bot *dto.AdminLineBot
	aru *admin.AdminReplyUsecase
}

// コンストラクタ
func NewAdminLineController(bot *dto.AdminLineBot, aru *admin.AdminReplyUsecase) *AdminLineController {
	return &AdminLineController{bot: bot, aru: aru}
}

func (alc *AdminLineController) Webhook(c *gin.Context) {
	events, err := alc.bot.ParseRequest(c.Request)
	if err != nil {
		conf.Log.Error("Failed to parse the request", zap.Any("err", err))
		return
	}
	conf.Log.Info("Successfully parse the request", zap.Any("events", events))
	for _, event := range events {
		if event.Source.Type == linebot.EventSourceTypeUser {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch event.Message.(type) {
				case *linebot.TextMessage:
					message := dto.NewTextMessage(event.ReplyToken, event.Message.(*linebot.TextMessage).Text)
					err = alc.aru.HandleTextMessage(message)
				default:
					message := dto.NewTextMessage(event.ReplyToken, "unknown")
					err = alc.aru.HandleTextMessage(message)
				}
			case linebot.EventTypePostback:
				message := dto.NewPostbackMessage(event.ReplyToken, event.Postback.Data, event.Postback.Params)
				err = alc.aru.HandlePostbackEvent(message)
			}
		}
	}
	if err != nil {
		conf.Log.Error("Failed to handle the request", zap.Any("err", err))
	}
}