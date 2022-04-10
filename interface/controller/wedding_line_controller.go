package controller

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"github.com/takakuwa-s/line-wedding-api/dto"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type WeddingLineController struct {
	bot *dto.WeddingLineBot
	wru *usecase.WeddingReplyUsecase
}

// コンストラクタ
func NewWeddingLineController(bot *dto.WeddingLineBot, wru *usecase.WeddingReplyUsecase) *WeddingLineController {
	return &WeddingLineController{bot: bot, wru: wru}
}

func (wlc *WeddingLineController) Webhook(c *gin.Context) {
	events, err := wlc.bot.ParseRequest(c.Request)
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
					file := entity.NewFile(event.Message.(*linebot.ImageMessage).ID, event.Source.UserID)
					message := dto.NewFileMessage(event.ReplyToken, file)
					err = wlc.wru.HandleFileEvent(message)
				case *linebot.VideoMessage:
					file := entity.NewFile(event.Message.(*linebot.ImageMessage).ID, event.Source.UserID)
					message := dto.NewFileMessage(event.ReplyToken, file)
					err = wlc.wru.HandleFileEvent(message)
				case *linebot.TextMessage:
					message := dto.NewTextMessage(event.ReplyToken, event.Message.(*linebot.TextMessage).Text, event.Source.UserID)
					err = wlc.wru.HandleTextMessage(message)
				default:
					message := dto.NewTextMessage(event.ReplyToken, "unknown", event.Source.UserID)
					err = wlc.wru.HandleTextMessage(message)
				}
			case linebot.EventTypePostback:
				message := dto.NewPostbackMessage(event.ReplyToken, event.Postback.Data, event.Source.UserID, event.Postback.Params)
				err = wlc.wru.HandlePostbackEvent(message)
			case linebot.EventTypeFollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID)
				err = wlc.wru.HandleFollowEvent(message)
			case linebot.EventTypeUnfollow:
				message := dto.NewFollowMessage(event.ReplyToken, event.Source.UserID)
				err = wlc.wru.HandleUnFollowEvent(message)
			}
		} else {
			message := dto.NewGroupMessage(event.ReplyToken)
			err = wlc.wru.HandleGroupEvent(message)
		}
	}
	if err != nil {
		conf.Log.Error("Failed to handle the request", zap.Any("err", err))
	}
}