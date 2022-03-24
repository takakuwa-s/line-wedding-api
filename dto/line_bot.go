package dto

import (
	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
)

type WeddingLineBot struct{
	*linebot.Client
}

type AdminLineBot struct{
	*linebot.Client
}

func NewWeddingLineBot() *WeddingLineBot {
	accessToken:= os.Getenv("WEDDING_BOT_ACCESS_TOKEN")
	channelSecret := os.Getenv("WEDDING_BOT_CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		conf.Log.Error("Failed to create lineBot instance", zap.Any("err", err), zap.Any("accessToken", accessToken), zap.Any("channelSecret", channelSecret))
	}
	return &WeddingLineBot{bot}
}

func NewAdminLineBot() *AdminLineBot {
	accessToken:= os.Getenv("ADMIN_BOT_ACCESS_TOKEN")
	channelSecret := os.Getenv("ADMIN_BOT_CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		conf.Log.Error("Failed to create lineBot instance", zap.Any("err", err), zap.Any("accessToken", accessToken), zap.Any("channelSecret", channelSecret))
	}
	return &AdminLineBot{bot}
}