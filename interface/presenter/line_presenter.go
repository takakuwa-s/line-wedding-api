package presenter

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type LinePresenter struct {
	bot *linebot.Client
}

// コンストラクタ
func NewLinePresenter() *LinePresenter {
	accessToken := os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
	}
	return &LinePresenter{bot: bot}
}

func (lp *LinePresenter) ReplyMessage(m dto.Message) {
	messages := linebot.NewTextMessage(m.Text)
	_, err := lp.bot.ReplyMessage(m.ReplyToken, messages).Do()
	if err != nil {
		return
	}
}