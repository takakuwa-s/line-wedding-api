package usecase

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type MessageHandler struct {
	ip ipresenter.IPresenter
}

// Newコンストラクタ
func NewMessageHandler(ip ipresenter.IPresenter) *MessageHandler {
	return &MessageHandler{ip:ip}
}

func (ml *MessageHandler) HandleTextMessage(m dto.Message) {
	ml.ip.ReplyMessage(m)
}
