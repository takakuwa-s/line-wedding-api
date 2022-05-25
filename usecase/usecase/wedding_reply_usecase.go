package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type WeddingReplyUsecase struct {
	mr  igateway.IMessageRepository
	lg  igateway.ILineGateway
	fg igateway.IFaceGateway
	ur  igateway.IUserRepository
	fr  igateway.IFileRepository
	br  igateway.IBinaryRepository
	apu *AdminPushUsecase
	cu  *CommonUtils
}

// Newコンストラクタ
func NewWeddingReplyUsecase(
	mr igateway.IMessageRepository,
	lg igateway.ILineGateway,
	fg igateway.IFaceGateway,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository,
	br igateway.IBinaryRepository,
	apu *AdminPushUsecase,
	cu *CommonUtils) *WeddingReplyUsecase {
	return &WeddingReplyUsecase{mr: mr, lg: lg, fg: fg, ur: ur, fr: fr, br: br, apu: apu, cu: cu}
}

func (wru *WeddingReplyUsecase) HandleImageEvent(m *dto.FileMessage) error {
	// Get the file binary
	content, err := wru.lg.GetFileContent(dto.WeddingBotType, m.File.Id)
	if err != nil {
		return err
	}

	// upload the file binary
	file, err := wru.br.SaveImageBinary(m.File, content)
	if err != nil {
		return err
	}
	faceRes, err := wru.fg.GetFaceAnalysis(file.ContentUrl)
	if err != nil {
		return err
	}
	wru.calcurateFaceScore(faceRes, file)

	// Save file data
	err = wru.fr.SaveFile(file)
	if err != nil {
		return err
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

	var profile *entity.User
	// follow the bot in the first time
	if user == nil {
		// Get the detail user profile
		profile, err = wru.lg.GetUserProfileById(m.SenderUserId, dto.WeddingBotType)
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
		profile = user
	}

	// Send notification to admin bot
	if err = wru.apu.SendFollowNotification(profile, user == nil); err != nil {
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
	return nil
}

func (wru *WeddingReplyUsecase) calcurateFaceScore(r []*dto.FaceResponse, f *entity.File) {
	if len(r) <= 0 || len(r) > 4 {
		f.FaceIds = []string{}
		f.FaceCount = 0
		f.FaceHappinessLevel = 0
		f.FacePhotoBeauty = 0
		f.FaceScore = 0
		return
	}
	faceCount := len(r)
	faceIds := make([]string, faceCount)
	var faceHappinessLevelSum float32
	var facePhotoBeautySum float32
	for i, f := range r {
		faceIds[i] = f.FaceId

		// calculate the face happiness level
		faceHappinessLevelSum += 10 * f.FaceAttributes.Smile
		faceHappinessLevelSum -= 3 * f.FaceAttributes.Emotion.Anger
		faceHappinessLevelSum -= 3 * f.FaceAttributes.Emotion.Contempt
		faceHappinessLevelSum -= 3 * f.FaceAttributes.Emotion.Disgust
		faceHappinessLevelSum -= f.FaceAttributes.Emotion.Fear
		faceHappinessLevelSum += 5 * f.FaceAttributes.Emotion.Happiness
		faceHappinessLevelSum += f.FaceAttributes.Emotion.Neutral
		faceHappinessLevelSum += 2 * f.FaceAttributes.Emotion.Surprise

		// calculate the face photo beauty
		facePhotoBeautySum += 10 * (1 - f.FaceAttributes.Blur.Value)
		facePhotoBeautySum += 10 + (1 - f.FaceAttributes.Noise.Value)
		switch f.FaceAttributes.Exposure.ExposureLevel {
			case "GoodExposure":
				facePhotoBeautySum += 5
			case "OverExposure":
				facePhotoBeautySum -= 2
			case "UnderExposure":
				facePhotoBeautySum -= 2
		}
		if f.FaceAttributes.Occlusion.ForeheadOccluded {
			facePhotoBeautySum -= 1
		}
		if f.FaceAttributes.Occlusion.EyeOccluded {
			facePhotoBeautySum -= 3
		}
		if f.FaceAttributes.Occlusion.MouthOccluded {
			facePhotoBeautySum -= 1
		}
	}
	f.FaceIds = faceIds
	f.FaceCount = faceCount
	f.FaceHappinessLevel = faceHappinessLevelSum / (0.95 * float32(faceCount))
	f.FacePhotoBeauty = facePhotoBeautySum / (0.95 * float32(faceCount))
	f.FaceScore = f.FaceHappinessLevel + f.FacePhotoBeauty
}