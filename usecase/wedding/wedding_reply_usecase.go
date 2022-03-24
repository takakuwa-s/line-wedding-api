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
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "image")
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
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "follow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), user.Name)
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

func (wru *WeddingReplyUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "group")
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (wru *WeddingReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := wru.mr.FindReplyMessage(dto.WeddingBotType, m.Text)
	if len(messages) == 0 {
		messages = wru.mr.FindMessageByKey(dto.WeddingBotType, "unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := wru.p.ReplyMessage(rm, dto.WeddingBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}
