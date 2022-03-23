package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)

type ReplyMessageUsecase struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
	lr igateway.ILineRepository
	ur igateway.IUserRepository
	fr igateway.IFileRepository
}

// Newコンストラクタ
func NewReplyMessageUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository) *ReplyMessageUsecase {
	return &ReplyMessageUsecase{p:p, mr:mr, lr:lr, ur:ur, fr:fr}
}

func (rmu *ReplyMessageUsecase) HandleFileEvent(m *dto.FileMessage) error {
	rmu.fr.SaveFile(m.File)
	messages := rmu.mr.FindImageMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := rmu.p.ReplyMessage(rm); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (rmu *ReplyMessageUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	user, err := rmu.lr.FindUserById(m.SenderUserId)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}
	user.CreatedAt = m.EventTime
	rmu.ur.SaveUser(user)
	messages := rmu.mr.FindFollowMessage(user.Name)
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err = rmu.p.ReplyMessage(rm); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (rmu *ReplyMessageUsecase) HandleUnFollowEvent(m *dto.FollowMessage) error {
	if err := rmu.ur.UpdateFollowStatusById(m.SenderUserId, false); err != nil {
		return fmt.Errorf("failed to update the follow status; err = %w", err)
	}
	return nil
}

func (rmu *ReplyMessageUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	// Do nothing
	return nil
}

func (rmu *ReplyMessageUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := rmu.mr.FindGroupMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := rmu.p.ReplyMessage(rm); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (rmu *ReplyMessageUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := rmu.mr.FindReplyMessage(m.Text)
	if len(messages) == 0 {
		messages = rmu.mr.FindReplyMessage("unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := rmu.p.ReplyMessage(rm); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}
