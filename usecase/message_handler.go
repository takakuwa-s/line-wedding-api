package usecase

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)

type MessageHandler struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
	lr igateway.ILineRepository
	ur igateway.IUserRepository
	fr igateway.IFileRepository
}

// Newコンストラクタ
func NewMessageHandler(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository) *MessageHandler {
	return &MessageHandler{p:p, mr:mr, lr:lr, ur:ur, fr:fr}
}

func (ml *MessageHandler) HandleFileEvent(m *dto.FileMessage) {
	ml.fr.SaveFile(m.File)
	messages := ml.mr.FindImageMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}

func (ml *MessageHandler) HandleFollowEvent(m *dto.FollowMessage) {
	user := ml.lr.FindUserById(m.SenderUserId)
	ml.ur.SaveUser(user)
	messages := ml.mr.FindFollowMessage(user.DisplayName)
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	ml.p.ReplyMessage(rm)
}

func (ml *MessageHandler) HandleUnFollowEvent(m *dto.FollowMessage) {
	ml.ur.UpdateFollowStatusById(m.SenderUserId, false)
}

func (ml *MessageHandler) HandlePostbackEvent(m *dto.PostbackMessage) {
	// Do nothing
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
