package admin

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)

type AdminReplyUsecase struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
}

// Newコンストラクタ
func NewAdminReplyUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository) *AdminReplyUsecase {
	return &AdminReplyUsecase{p:p, mr:mr}
}

func (aru *AdminReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := aru.mr.FindReplyMessage(m.Text)
	if len(messages) == 0 {
		messages = aru.mr.FindReplyMessage("unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := aru.p.ReplyMessage(rm, dto.AdminBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}
