package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"go.uber.org/zap"
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
	} else if m.Text == "写真を削除" {
		messages = wru.handleImageDeleteRequest(messages, m.SenderUserId)
	}
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	var messages []map[string]interface{}
	switch m.Data["action"].(string) {
	case "image_delete":
		fileId := m.Data["fileId"].(string)
		lineFileId := m.Data["lineFileId"].(string)
		if err := wru.br.DeleteBinary(fileId); err != nil {
			conf.Log.Error("Failed to delete image binary", zap.Any("err", err))
			messages = wru.mr.FindMessageByKey(dto.WeddingBotType, "delete_image_error")
			break
		}
		if err := wru.fr.DeleteFile(lineFileId, m.SenderUserId); err != nil {
			conf.Log.Error("Failed to delete the file metadata", zap.Any("err", err))
		}
		messages = wru.mr.FindMessageByKey(dto.WeddingBotType, "delete_image_success")
	case "image_check":
		messages = wru.mr.FindMessageByKey(dto.WeddingBotType, "image_check")
		messages[0]["originalContentUrl"] = fmt.Sprintf(messages[0]["originalContentUrl"].(string), m.Data["id"].(string))
		messages[0]["previewImageUrl"] = fmt.Sprintf(messages[0]["previewImageUrl"].(string), m.Data["id"].(string))
	}
	return wru.cu.SendReplyMessage(m.ReplyToken, messages, dto.WeddingBotType)
}

func (wru *WeddingReplyUsecase) handleImageDeleteRequest(m []map[string]interface{}, userId string) []map[string]interface{} {
	files, err := wru.fr.FindByCreaterAndIsDeleted(userId, false)
	if err != nil {

	}
	if len(files) == 0 {
		return wru.mr.FindMessageByKey(dto.WeddingBotType, "not_found_image")
	}
	cols := m[0]["template"].(map[string]interface{})["columns"].([]interface{})
	byteCol, err := json.Marshal(cols[0].(map[string]interface{}))
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal column; err = %v", err))
	}
	for _, file := range files {
		var obj interface{}
  	if err := json.Unmarshal(byteCol, &obj); err != nil {
			panic(fmt.Sprintf("Failed to unmarshal column; err = %v", err))
		}
		c := obj.(map[string]interface{})
		c["thumbnailImageUrl"] = fmt.Sprintf(c["thumbnailImageUrl"].(string), file.FileId)
		c["text"] = fmt.Sprintf(c["text"].(string), file.Name)
		actions := c["actions"].([]interface{})
		actions[0].(map[string]interface{})["data"] = fmt.Sprintf(actions[0].(map[string]interface{})["data"].(string), file.FileId)
		actions[1].(map[string]interface{})["data"] = fmt.Sprintf(actions[1].(map[string]interface{})["data"].(string), file.LineFileId, file.FileId)
		cols = append(cols, c)
	}
	m[0]["template"].(map[string]interface{})["columns"] = cols[1:]
	return m
}