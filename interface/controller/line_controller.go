package controller

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type LineController struct {
	bot *linebot.Client
	rmu  *usecase.ReplyMessageUsecase
}

// コンストラクタ
func NewLineController(bot *linebot.Client, rmu *usecase.ReplyMessageUsecase) *LineController {
	return &LineController{bot: bot, rmu: rmu}
}

func (lw *LineController) Webhook(c *gin.Context) {
	events, err := lw.bot.ParseRequest(c.Request)
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
				case *linebot.ImageMessage:
					file := entity.NewFile(event.Message.(*linebot.ImageMessage).ID, event.Source.UserID, event.Timestamp)
					message := dto.NewFileMessage(event.ReplyToken, file)
					err = lw.rmu.HandleFileEvent(message)
				case *linebot.VideoMessage:
					file := entity.NewFile(event.Message.(*linebot.ImageMessage).ID, event.Source.UserID, event.Timestamp)
					message := dto.NewFileMessage(event.ReplyToken, file)
					err = lw.rmu.HandleFileEvent(message)
				case *linebot.TextMessage:
					message := dto.NewTextMessage(event.ReplyToken, event.Message.(*linebot.TextMessage).Text)
					err = lw.rmu.HandleTextMessage(message)
				default:
					message := dto.NewTextMessage(event.ReplyToken, "unknown")
					err = lw.rmu.HandleTextMessage(message)
				}
			case linebot.EventTypeFollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID, event.Timestamp)
				err = lw.rmu.HandleFollowEvent(message)
			case linebot.EventTypeUnfollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID, event.Timestamp)
				err = lw.rmu.HandleUnFollowEvent(message)
			case linebot.EventTypePostback:
				message := dto.NewPostbackMessage(event.ReplyToken, event.Postback.Data, event.Postback.Params)
				err = lw.rmu.HandlePostbackEvent(message)
			}
		} else {
			message := dto.NewGroupMessage(event.ReplyToken)
			err = lw.rmu.HandleGroupEvent(message)
		}
	}
	if err != nil {
		conf.Log.Error("Failed to handle the request", zap.Any("err", err))
	}
}