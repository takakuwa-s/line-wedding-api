package dto

import (
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LineBot struct {
	*linebot.Client
}

func NewLineBot() *LineBot {
	accessToken := os.Getenv("LINE_BOT_ACCESS_TOKEN")
	channelSecret := os.Getenv("LINE_BOT_CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the wedding lineBot instance; err = %v", err))
	}
	return &LineBot{bot}
}
