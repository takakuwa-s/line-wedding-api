package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"go.uber.org/zap"
)

type LineReplyUsecase struct {
	mr  igateway.IMessageRepository
	lg  igateway.ILineGateway
	fg  igateway.IFaceGateway
	ur  igateway.IUserRepository
	fr  igateway.IFileRepository
	br  igateway.IBinaryRepository
	lpu *LinePushUsecase
	p   ipresenter.IPresenter
}

// Newコンストラクタ
func NewLineReplyUsecase(
	mr igateway.IMessageRepository,
	lg igateway.ILineGateway,
	fg igateway.IFaceGateway,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository,
	br igateway.IBinaryRepository,
	lpu *LinePushUsecase,
	p   ipresenter.IPresenter) *LineReplyUsecase {
	return &LineReplyUsecase{mr: mr, lg: lg, fg: fg, ur: ur, fr: fr, br: br, lpu: lpu, p: p}
}

func (lru *LineReplyUsecase) HandleImageEvent(m *dto.FileMessage) error {
	// Get the file binary
	content, err := lru.lg.GetFileContent(m.File.Id)
	if err != nil {
		return err
	}

	// upload the file binary
	file, err := lru.br.SaveImageBinary(m.File, content)
	if err != nil {
		return err
	}
	faceRes, err := lru.fg.GetFaceAnalysis(file.ContentUrl)
	if err != nil {
		return err
	}
	lru.calcurateFaceScore(faceRes, file)

	// Save file data
	err = lru.fr.SaveFile(file)
	if err != nil {
		return err
	}

	// Reply message
	messages := lru.mr.FindMessageByKey("image")
	return lru.sendReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := lru.ur.FindById(m.SenderUserId)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}

	var profile *entity.User
	// follow the bot in the first time
	if user == nil {
		// Get the detail user profile
		profile, err = lru.lg.GetUserProfileById(m.SenderUserId)
		if err != nil {
			return fmt.Errorf("failed to find the user; err = %w", err)
		}

		// Save users
		if err = lru.ur.SaveUser(profile); err != nil {
			return fmt.Errorf("failed to save a user; err = %w", err)
		}
	} else {
		// update user status
		if err := lru.ur.UpdateFollowById(m.SenderUserId, true); err != nil {
			return fmt.Errorf("failed to update the follow status; err = %w", err)
		}
		profile = user
	}

	// Send notification to admin bot
	if err = lru.lpu.SendFollowNotification(profile, user == nil); err != nil {
		return fmt.Errorf("failed to send notification to admin bot; err = %w", err)
	}

	// Return message
	messages := lru.mr.FindMessageByKey("follow")
	messages[0]["text"] = fmt.Sprintf(messages[0]["text"].(string), user.LineName)
	return lru.sendReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleUnFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := lru.ur.FindById(m.SenderUserId)
	if err != nil {
		return fmt.Errorf("failed to find the user; err = %w", err)
	}
	if user == nil {
		return fmt.Errorf("not found the user")
	}

	// update user status
	if err := lru.ur.UpdateFollowById(m.SenderUserId, false); err != nil {
		return fmt.Errorf("failed to update the follow status; err = %w", err)
	}

	// Send notification to admin bot
	if err := lru.lpu.SendUnFollowNotification(user); err != nil {
		return fmt.Errorf("failed to send notification to admin bot; err = %w", err)
	}

	return nil
}

func (lru *LineReplyUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := lru.mr.FindMessageByKey("group")
	return lru.sendReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	var messages []map[string]interface{}
	switch m.Text {
	case "招待状送信内容確認":
		if lru.checkAdminRole(m.SenderUserId) {
			messages = lru.mr.FindMessageByKey("invitation")
		}
	case "前日メッセージ送信内容確認":
		if lru.checkAdminRole(m.SenderUserId) {
			messages = lru.mr.FindMessageByKey("reminder")
		}
	default:
		messages = lru.mr.FindReplyMessage(m.Text)
	}
	if len(messages) == 0 {
		messages = lru.mr.FindMessageByKey("unknown")
	}
	return lru.sendReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	return nil
}

func (lru *LineReplyUsecase) sendReplyMessage(
	token string,
	m []map[string]interface{}) error {
	rm := dto.NewReplyMessage(token, m)
	if err := lru.p.ReplyMessage(rm); err != nil {
		return err
	}
	return nil
}

func (lru *LineReplyUsecase) checkAdminRole(userId string) bool {
	user, err := lru.ur.FindById(userId)
	if err != nil {
		conf.Log.Error("failed to find the user", zap.Error(err))
		return false
	}
	return user != nil && user.IsAdmin
}

func (lru *LineReplyUsecase) calcurateFaceScore(r []*dto.FaceResponse, f *entity.File) {
	if len(r) <= 0 || len(r) > 10 {
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
	var hasMale bool
	var hasFemale bool
	var hasYoung bool
	var hasElderly bool
	for i, f := range r {
		faceIds[i] = f.FaceId

		// calculate the face happiness level (max: 40)
		faceHappinessLevelSum += 20 * f.FaceAttributes.Smile
		faceHappinessLevelSum -= 20 * f.FaceAttributes.Emotion.Anger
		faceHappinessLevelSum -= 10 * f.FaceAttributes.Emotion.Contempt
		faceHappinessLevelSum -= 15 * f.FaceAttributes.Emotion.Disgust
		faceHappinessLevelSum -= 5 * f.FaceAttributes.Emotion.Fear
		faceHappinessLevelSum += 20 * f.FaceAttributes.Emotion.Happiness
		faceHappinessLevelSum += 1 * f.FaceAttributes.Emotion.Neutral
		faceHappinessLevelSum += 5 * f.FaceAttributes.Emotion.Surprise

		// calculate the face photo beauty (max: 30)
		facePhotoBeautySum += 10 * (1 - f.FaceAttributes.Blur.Value)
		facePhotoBeautySum += 10 * (1 - f.FaceAttributes.Noise.Value)
		switch f.FaceAttributes.Exposure.ExposureLevel {
			case "GoodExposure":
				facePhotoBeautySum += 10
			case "OverExposure":
				facePhotoBeautySum -= 5
			case "UnderExposure":
				facePhotoBeautySum -= 5
		}
		if f.FaceAttributes.Occlusion.ForeheadOccluded {
			facePhotoBeautySum -= 2
		}
		if f.FaceAttributes.Occlusion.EyeOccluded {
			facePhotoBeautySum -= 4
		}
		if f.FaceAttributes.Occlusion.MouthOccluded {
			facePhotoBeautySum -= 2
		}

		// For bonus
		if f.FaceAttributes.Gender == "male" {
			hasMale = true
		}
		if f.FaceAttributes.Gender == "female" {
			hasFemale = true
		}
		if f.FaceAttributes.Age < 10 {
			hasYoung = true
		}
		if f.FaceAttributes.Age > 50 {
			hasElderly = true
		}
	}
	// calculate the face count bonus point (max: 20)
	bonusPoint := 2 * float32(faceCount)
	if hasMale && hasFemale {
		bonusPoint += 4
	}
	if hasYoung {
		bonusPoint += 3
	}
	if hasElderly {
		bonusPoint += 3
	}

	f.FaceCount = faceCount
	f.FaceHappinessLevel = faceHappinessLevelSum / float32(faceCount)
	f.FacePhotoBeauty = facePhotoBeautySum / float32(faceCount)
	f.FaceScore = f.FaceHappinessLevel + f.FacePhotoBeauty + bonusPoint
}