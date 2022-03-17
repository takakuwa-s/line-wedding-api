package presenter

import (
	conf "github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type LinePresenter struct {
	bot *linebot.Client
}

// コンストラクタ
func NewLinePresenter(bot *linebot.Client) *LinePresenter {
	return &LinePresenter{bot: bot}
}

func (lp *LinePresenter) ReplyMessage(m dto.Message) {
	messages := linebot.NewTextMessage(m.Text)
	if _, err := lp.bot.ReplyMessage(m.ReplyToken, messages).Do(); err != nil {
		conf.Log.Error("Failed to send the reply message", zap.Any("err", err), zap.Any("messages", messages))
	}
}