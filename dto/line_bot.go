package dto

import (
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type WeddingLineBot struct {
	*linebot.Client
}

type AdminLineBot struct {
	*linebot.Client
}

func NewWeddingLineBot() *WeddingLineBot {
	accessToken := os.Getenv("WEDDING_BOT_ACCESS_TOKEN")
	channelSecret := os.Getenv("WEDDING_BOT_CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the wedding lineBot instance; err = %v", err))
	}
	return &WeddingLineBot{bot}
}

func NewAdminLineBot() *AdminLineBot {
	accessToken := os.Getenv("ADMIN_BOT_ACCESS_TOKEN")
	channelSecret := os.Getenv("ADMIN_BOT_CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		panic(fmt.Sprintf("Failed to create the admin lineBot instance; err = %v", err))
	}
	return &AdminLineBot{bot}
}
