package usecase

import (
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
)

type AdminReplyUsecase struct {
	mr  igateway.IMessageRepository
	lr  igateway.ILineRepository
	wpu *WeddingPushUsecase
	cu  *CommonUtils
}

// Newコンストラクタ
func NewAdminReplyUsecase(
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	wpu *WeddingPushUsecase,
	cu *CommonUtils) *AdminReplyUsecase {
	return &AdminReplyUsecase{mr: mr, lr: lr, wpu: wpu, cu:cu}
}

func (aru *AdminReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	var messages []map[string]interface{}
	switch m.Data["action"].(string) {
	case "invitation":
		if m.Data["confirm"].(bool) {
			if err := aru.wpu.PublishInvitation(); err != nil {
				messages = aru.mr.FindMessageByKey(dto.AdminBotType, "invitation_error")
				messages[1]["text"] = fmt.Sprintf(messages[1]["text"].(string), err)
			} else {
				messages = aru.mr.FindMessageByKey(dto.AdminBotType, "invitation_submit")
			}
		} else {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "postback_cancel")
		}
	case "reminder":
		if m.Data["confirm"].(bool) {
			if err := aru.wpu.PublishReminder(); err != nil {
				messages = aru.mr.FindMessageByKey(dto.AdminBotType, "reminder_error")
				messages[1]["text"] = fmt.Sprintf(messages[1]["text"].(string), err)
			} else {
				messages = aru.mr.FindMessageByKey(dto.AdminBotType, "reminder_submit")
			}
		} else {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "postback_cancel")
		}
	case "slideshow":
		if m.Data["confirm"].(bool) {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "slideshow_submit")
		} else {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "postback_cancel")
		}
	}
	return aru.cu.SendReplyMessage(m.ReplyToken, messages, dto.AdminBotType)
}

func (aru *AdminReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := aru.mr.FindReplyMessage(dto.AdminBotType, m.Text)
	var err error
	switch m.Text {
	case "送信済みpushメッセージ数":
		messages, err = aru.handleQuotaCheck(messages)
	case "招待状を確認":
		messages = aru.mr.FindMessageByKey(dto.WeddingBotType, "invitation")
	case "前日メッセージを確認":
		messages = aru.mr.FindMessageByKey(dto.WeddingBotType, "reminder")
	case "スライドショーを確認":
	}
	if err != nil {
		return err
	}
	if len(messages) == 0 {
		messages = aru.mr.FindMessageByKey(dto.AdminBotType, "unknown")
	}
	return aru.cu.SendReplyMessage(m.ReplyToken, messages, dto.AdminBotType)
}

func (aru *AdminReplyUsecase) handleQuotaCheck(m []map[string]interface{}) ([]map[string]interface{}, error) {
	weddingComsuption, err := aru.lr.GetQuotaComsuption(dto.WeddingBotType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the wedding quota comsuption; err = %w", err)
	}
	adminComsuption, err := aru.lr.GetQuotaComsuption(dto.AdminBotType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the admin quota comsuption; err = %w", err)
	}
	m[0]["text"] = fmt.Sprintf(m[0]["text"].(string), int(weddingComsuption))
	m[1]["text"] = fmt.Sprintf(m[1]["text"].(string), int(adminComsuption))
	return m, nil
}