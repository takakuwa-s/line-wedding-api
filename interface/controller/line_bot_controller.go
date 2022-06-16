package controller

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type LineBotController struct {
	lb  *dto.LineBot
	lru *usecase.LineReplyUsecase
}

// コンストラクタ
func NewLineBotController(lb *dto.LineBot, lru *usecase.LineReplyUsecase) *LineBotController {
	return &LineBotController{lb: lb, lru: lru}
}

func (lbc *LineBotController) Webhook(c *gin.Context) {
	bot, err := lbc.lb.GetClient()
	if err != nil {
		conf.Log.Error("Failed to get the line bot instance", zap.Error(err))
		return
	}
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		conf.Log.Error("Failed to parse the request", zap.Error(err))
		return
	}
	conf.Log.Info("Successfully parse the request", zap.Any("events", events))
	for _, event := range events {
		if event.Source.Type == linebot.EventSourceTypeUser {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch event.Message.(type) {
				case *linebot.ImageMessage:
					imageMessage := event.Message.(*linebot.ImageMessage)
					file := entity.NewFile(imageMessage.ID, event.Source.UserID)
					var imageSet *entity.ImageSet
					if imageMessage.ImageSet != nil {
						imageSet = entity.NewImageSet(imageMessage.ImageSet.ID, imageMessage.ImageSet.Total)
					}
					message := dto.NewFileMessage(event.ReplyToken, file, imageSet)
					err = lbc.lru.HandleImageEvent(message)
				case *linebot.TextMessage:
					message := dto.NewTextMessage(event.ReplyToken, event.Message.(*linebot.TextMessage).Text, event.Source.UserID)
					err = lbc.lru.HandleTextMessage(message)
				default:
					message := dto.NewTextMessage(event.ReplyToken, "unknown", event.Source.UserID)
					err = lbc.lru.HandleTextMessage(message)
				}
			case linebot.EventTypePostback:
				message := dto.NewPostbackMessage(event.ReplyToken, event.Postback.Data, event.Source.UserID, event.Postback.Params)
				err = lbc.lru.HandlePostbackEvent(message)
			case linebot.EventTypeFollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID)
				err = lbc.lru.HandleFollowEvent(message)
			case linebot.EventTypeUnfollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID)
				err = lbc.lru.HandleUnFollowEvent(message)
			}
		} else {
			message := dto.NewGroupMessage(event.ReplyToken)
			err = lbc.lru.HandleGroupEvent(message)
		}
		if err != nil {
			conf.Log.Error("Failed to handle the request", zap.Error(err))
			lbc.lru.HandleError(event.ReplyToken)
		}
	}
}
