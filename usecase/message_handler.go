package usecase

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)

type MessageHandler struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
}

// Newコンストラクタ
func NewMessageHandler(p ipresenter.IPresenter, mr igateway.IMessageRepository) *MessageHandler {
	return &MessageHandler{p:p, mr:mr}
}

func (ml *MessageHandler) HandleTextMessage(m dto.RequestMessage) {
	messages := ml.mr.FindReplyMessage(m.Text)
	if len(messages) == 0 {
		messages = ml.mr.FindReplyMessage("unknown")
	}
	rm := dto.ReplyMessage{
		ReplyToken: m.ReplyToken,
		Messages: messages,
	}
	ml.p.ReplyMessage(rm)
}
