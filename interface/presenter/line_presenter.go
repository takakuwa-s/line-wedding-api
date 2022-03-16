package presenter

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"os"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
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
		logger.Error("Failed to create lineBot instance", zap.Any("err", err))
	}
	return &LinePresenter{bot: bot}
}

func (lp *LinePresenter) ReplyMessage(m dto.Message) {
	messages := linebot.NewTextMessage(m.Text)
	_, err := lp.bot.ReplyMessage(m.ReplyToken, messages).Do()
	if err != nil {
		logger.Error("Failed to send the reply message", zap.Any("err", err), zap.Any("messages", messages))
		return
	}
}