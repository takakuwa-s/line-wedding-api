package wedding

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
)

type WeddingReplyUsecase struct {
	p ipresenter.IPresenter
	mr igateway.IMessageRepository
	lr igateway.ILineRepository
	ur igateway.IUserRepository
	fr igateway.IFileRepository
}

// Newコンストラクタ
func NewWeddingReplyUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository) *WeddingReplyUsecase {
	return &WeddingReplyUsecase{p:p, mr:mr, lr:lr, ur:ur, fr:fr}
}

func (wru *WeddingReplyUsecase) HandleFileEvent(m *dto.FileMessage) error {
	wru.fr.SaveFile(m.File)
	messages := wru.mr.FindImageMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (wru *WeddingReplyUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	user, err := wru.lr.FindUserById(m.SenderUserId, dto.WeddingBotType)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}
	user.CreatedAt = m.EventTime
	wru.ur.SaveUser(user)
	messages := wru.mr.FindFollowMessage(user.Name)
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err = wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (wru *WeddingReplyUsecase) HandleUnFollowEvent(m *dto.FollowMessage) error {
	if err := wru.ur.UpdateFollowStatusById(m.SenderUserId, false); err != nil {
		return fmt.Errorf("failed to update the follow status; err = %w", err)
	}
	return nil
}

func (wru *WeddingReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	// Do nothing
	return nil
}

func (wru *WeddingReplyUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := wru.mr.FindGroupMessage()
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (wru *WeddingReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := wru.mr.FindReplyMessage(m.Text)
	if len(messages) == 0 {
		messages = wru.mr.FindReplyMessage("unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}
