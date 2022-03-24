package admin

import (
	"encoding/json"
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/wedding"
)

type AdminReplyUsecase struct {
	p   ipresenter.IPresenter
	mr  igateway.IMessageRepository
	lr  igateway.ILineRepository
	wpu *wedding.WeddingPushUsecase
}

// Newコンストラクタ
func NewAdminReplyUsecase(
	p ipresenter.IPresenter,
	mr igateway.IMessageRepository,
	lr igateway.ILineRepository,
	wpu *wedding.WeddingPushUsecase) *AdminReplyUsecase {
	return &AdminReplyUsecase{p: p, mr: mr, lr: lr, wpu: wpu}
}

func (aru *AdminReplyUsecase) HandlePostbackEvent(m *dto.PostbackMessage) error {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(m.Data), &data); err != nil {
		return fmt.Errorf("failed to convert postback data to map object; err = %w", err)
	}
	var messages []map[string]interface{}
	switch data["action"].(string) {
	case "reminder":
		if data["confirm"].(bool) {
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
		if data["confirm"].(bool) {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "slideshow_submit")
		} else {
			messages = aru.mr.FindMessageByKey(dto.AdminBotType, "postback_cancel")
		}
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := aru.p.ReplyMessage(rm, dto.AdminBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
}

func (aru *AdminReplyUsecase) HandleTextMessage(m *dto.TextMessage) error {
	messages := aru.mr.FindReplyMessage(dto.AdminBotType, m.Text)
	var err error
	switch m.Text {
	case "残pushメッセージ数":
		messages, err = aru.handleQuotaCheck(messages)
	case "前日メッセージを確認":
		messages = aru.mr.FindMessageByKey(dto.WeddingBotType, "reminder")
	case "スライドショーを確認":
	}
	if err != nil {
		return fmt.Errorf("; err = %w", err)
	}
	if len(messages) == 0 {
		messages = aru.mr.FindMessageByKey(dto.AdminBotType, "unknown")
	}
	rm := dto.NewReplyMessage(m.ReplyToken, messages)
	if err := aru.p.ReplyMessage(rm, dto.AdminBotType); err != nil {
		return fmt.Errorf("failed to send the reply message; err = %w", err)
	}
	return nil
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
