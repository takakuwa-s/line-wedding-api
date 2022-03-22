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

func (ml *MessageHandler) HandleImageEvent(m *dto.FileMessage) {
	messages := ml.mr.FindImageMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}

func (ml *MessageHandler) HandleFollowEvent(m *dto.FollowMessage) {
	messages := ml.mr.FindFollowMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}

func (ml *MessageHandler) HandleGroupEvent(m *dto.GroupMessage) {
	messages := ml.mr.FindGroupMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}

func (ml *MessageHandler) HandleTextMessage(m *dto.TextMessage) {
	messages := ml.mr.FindReplyMessage(m.Text)
	if len(messages) == 0 {
		messages = ml.mr.FindReplyMessage("unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}
