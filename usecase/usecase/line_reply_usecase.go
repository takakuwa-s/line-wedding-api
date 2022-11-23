package usecase

import (
	"fmt"
	"os"

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
	ur  igateway.IUserRepository
	fr  igateway.IFileRepository
	isr igateway.IImageSetRepository
	bpg igateway.IBackgroundProcessGateway
	lpu *LinePushUsecase
	su  *SlideShowUsecase
	p   ipresenter.IPresenter
}

// Newコンストラクタ
func NewLineReplyUsecase(
	mr igateway.IMessageRepository,
	lg igateway.ILineGateway,
	ur igateway.IUserRepository,
	fr igateway.IFileRepository,
	isr igateway.IImageSetRepository,
	bpg igateway.IBackgroundProcessGateway,
	lpu *LinePushUsecase,
	su *SlideShowUsecase,
	p ipresenter.IPresenter) *LineReplyUsecase {
	return &LineReplyUsecase{mr: mr, lg: lg, ur: ur, fr: fr, isr: isr, bpg: bpg, lpu: lpu, su: su, p: p}
}

func (lru *LineReplyUsecase) HandleImageEvent(m *dto.FileMessage) error {
	// Save file data
	err := lru.fr.SaveFile(m.File)
	if err != nil {
		return err
	}

	if m.ImageSet == nil {
		// Start file uploading
		err = lru.bpg.StartUploadingFiles([]string{m.File.Id})
		if err != nil {
			lru.fr.DeleteFileById(m.File.Id)
			return err
		}
	} else {
		// Save image set
		imageSet, err := lru.isr.AppendFileIdByImageSet(m.ImageSet, m.File.Id)
		if err != nil {
			return err
		}
		if len(imageSet.FileIds) < imageSet.Total {
			return nil
		}
		err = lru.isr.DeleteById(m.ImageSet.Id)
		if err != nil {
			conf.Log.Error("failed to delete the image set", zap.Error(err))
		}
		// Start file uploading
		err = lru.bpg.StartUploadingFiles(imageSet.FileIds)
		if err != nil {
			lru.fr.DeleteFileByIds(imageSet.FileIds)
			return err
		}
	}

	// Reply message
	messages := lru.mr.FindMessageByKey("image")
	return lru.p.ReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleVideoEvent(m *dto.FileMessage) error {
	if m.File.Duration > 5*60*1000 {
		messages := lru.mr.FindMessageByKey("video_error")
		return lru.p.ReplyMessage(m.ReplyToken, messages)
	}

	// Save file data
	err := lru.fr.SaveFile(m.File)
	if err != nil {
		return err
	}

	// Start file uploading
	err = lru.bpg.StartUploadingFiles([]string{m.File.Id})
	if err != nil {
		lru.fr.DeleteFileById(m.File.Id)
		return err
	}

	// Reply message
	messages := lru.mr.FindMessageByKey("video")
	return lru.p.ReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := lru.ur.FindById(m.SenderUserId)
	if err != nil {
		return err
	}

	var profile *entity.User
	// follow the bot in the first time
	if user == nil {
		// Get the detail user profile
		profile, err = lru.lg.GetUserProfileById(m.SenderUserId)
		if err != nil {
			return err
		}

		// Save users
		if err = lru.ur.SaveUser(profile); err != nil {
			return err
		}
	} else {
		// update user status
		if err := lru.ur.UpdateBoolFieldById(m.SenderUserId, "Follow", true); err != nil {
			return err
		}
		profile = user
	}

	// Send notification to admin bot
	if err = lru.lpu.SendFollowNotification(profile, user == nil); err != nil {
		return fmt.Errorf("failed to send notification to admin user; err = %w", err)
	}

	// Return message
	messages := lru.mr.FindMessageByKey("follow")
	liffUrl := os.Getenv("LIFF_URL")
	messages[2]["text"] = fmt.Sprintf(messages[2]["text"].(string), liffUrl)
	return lru.p.ReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleUnFollowEvent(m *dto.FollowMessage) error {
	// Get user
	user, err := lru.ur.FindById(m.SenderUserId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("not found the user")
	}

	// update user status
	if err := lru.ur.UpdateBoolFieldById(m.SenderUserId, "Follow", false); err != nil {
		return err
	}

	// Send notification to admin bot
	if err := lru.lpu.SendUnFollowNotification(user); err != nil {
		return fmt.Errorf("failed to send notification to admin user; err = %w", err)
	}

	return nil
}

func (lru *LineReplyUsecase) HandleGroupEvent(m *dto.GroupMessage) error {
	messages := lru.mr.FindMessageByKey("group")
	return lru.p.ReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	var messages []map[string]interface{}
	switch m.Text {
	case "前日メッセージ送信内容確認":
		if lru.checkAdminRole(m.SenderUserId) {
			messages = lru.mr.FindMessageByKey("reminder")
		}
	case "スライドショー確認":
		if lru.checkAdminRole(m.SenderUserId) {
			var err error
			messages, err = lru.su.CreateSlideshowMessage()
			if err != nil {
				return err
			}
		}
	default:
		messages = lru.mr.FindReplyMessage(m.Text)
	}
	if len(messages) == 0 {
		messages = lru.mr.FindMessageByKey("unknown")
	}
	return lru.p.ReplyMessage(m.ReplyToken, messages)
}

func (lru *LineReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	return nil
}

func (lru *LineReplyUsecase) HandleError(token string) {
	messages := lru.mr.FindMessageByKey("error")
	if err := lru.p.ReplyMessage(token, messages); err != nil {
		conf.Log.Error("failed to send error reply message", zap.Error(err))
	}
}

func (lru *LineReplyUsecase) checkAdminRole(userId string) bool {
	user, err := lru.ur.FindById(userId)
	if err != nil {
		return false
	}
	return user != nil && user.IsAdmin
}
