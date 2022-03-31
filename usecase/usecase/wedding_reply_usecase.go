package usecase

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
	apu *AdminPushUsecase
	cu *CommonUsecase
}

// Newコンストラクタ
func NewWeddingReplyUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository,
	apu *AdminPushUsecase,
	cu *CommonUsecase) *WeddingReplyUsecase {
	return &WeddingReplyUsecase{p:p, mr:mr, lr:lr, ur:ur, fr:fr, apu:apu, cu:cu}
}

func (wru *WeddingReplyUsecase) HandleFileEvent(m *dto.FileMessage) error {
	wru.fr.SaveFile(m.File)
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "image")
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	// Get user profile
	user, err := wru.lr.FindUserById(m.SenderUserId, dto.WeddingBotType)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}

	// Save user
	user.CreatedAt = m.EventTime
	if err = wru.ur.SaveUser(user); err != nil {
		return fmt.Errorf("failed to save a user; err = %w", err)
	}

	// Send notification to admin bot
	if err = wru.apu.SendFollowNotification(user); err != nil {
		return fmt.Errorf("failed to send notification to admin bot; err = %w", err)
	}
	
	// Return message
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "follow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), user.Name)
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) HandleUnFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := wru.ur.FindById(m.SenderUserId)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}

	// update user status
	if err := wru.ur.UpdateFollowStatusById(m.SenderUserId, false); err != nil {
		return fmt.Errorf("failed to update the follow status; err = %w", err)
	}

	// Send notification to admin bot
	if err := wru.apu.SendUnFollowNotification(user); err != nil {
		return fmt.Errorf("failed to send notification to admin bot; err = %w", err)
	}

	return nil
}

func (wru *WeddingReplyUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "group")
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := wru.mr.FindReplyMessage(dto.WeddingBotType, m.Text)
	if len(messages) == 0 {
		messages = wru.mr.FindMessageByKey(dto.WeddingBotType, "unknown")
	}
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}
