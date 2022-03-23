package conf

import (
	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

var (
	Log, _ = zap.NewProduction()
)

func NewLineBot() *linebot.Client {
	accessToken:= os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		Log.Error("Failed to create lineBot instance", zap.Any("err", err), zap.Any("accessToken", accessToken), zap.Any("channelSecret", channelSecret))
	}
	return bot
}