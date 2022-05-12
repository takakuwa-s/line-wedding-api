package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type WeddingReplyUsecase struct {
	mr  igateway.IMessageRepository
	lr  igateway.ILineRepository
	ur  igateway.IUserRepository
	fr  igateway.IFileRepository
	br  igateway.IBinaryRepository
	apu *AdminPushUsecase
	cu  *CommonUtils
}

// Newコンストラクタ
func NewWeddingReplyUsecase(
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository,
	br igateway.IBinaryRepository,
	apu *AdminPushUsecase,
	cu *CommonUtils) *WeddingReplyUsecase {
	return &WeddingReplyUsecase{mr: mr, lr: lr, ur: ur, fr: fr, br: br, apu: apu, cu: cu}
}

func (wru *WeddingReplyUsecase) HandleFileEvent(m *dto.FileMessage) error {
	// Get the file binary
	content, err := wru.lr.GetFileContent(dto.WeddingBotType, m.File.LineFileId)
	if err != nil {
		return fmt.Errorf("failed to download the file; err = %w", err)
	}

	// upload the file binary
	file, err := wru.br.SaveBinary(m.File, content)
	if err != nil {
		return fmt.Errorf("failed to upload the file; err = %w", err)
	}

	// Save file data
	err = wru.fr.SaveFile(file)
	if err != nil {
		return fmt.Errorf("failed to save the file; err = %w", err)
	}

	// Reply message
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "image")
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := wru.ur.FindById(m.SenderUserId)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}

	// follow the bot in the first time
	if user == nil {
		// Get the detail user profile
		profile, err := wru.lr.GetUserProfileById(m.SenderUserId, dto.WeddingBotType)
		if err != nil {
			return fmt.Errorf("failed to find the user; err = %w", err)
		}

		// Save users
		if err = wru.ur.SaveUser(profile); err != nil {
			return fmt.Errorf("failed to save a user; err = %w", err)
		}
	} else {
		// update user status
		if err := wru.ur.UpdateFollowStatusById(m.SenderUserId, true); err != nil {
			return fmt.Errorf("failed to update the follow status; err = %w", err)
		}
	}

	// Send notification to admin bot
	if err = wru.apu.SendFollowNotification(user, user == nil); err != nil {
		return fmt.Errorf("failed to send notification to admin bot; err = %w", err)
	}

	// Return message
	messages := wru.mr.FindMessageByKey(dto.WeddingBotType, "follow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), user.LineName)
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

func (wru *WeddingReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	return fmt.Errorf("not implemented")
}