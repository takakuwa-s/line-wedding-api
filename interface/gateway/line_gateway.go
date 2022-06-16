package gateway

import (
	"fmt"
	"io"

	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type LineGateway struct {
	lb *dto.LineBot
}

// Newコンストラクタ
func NewLineGateway(lb *dto.LineBot) *LineGateway {
	return &LineGateway{lb: lb}
}

func (lg *LineGateway) GetUserProfileById(id string) (*entity.User, error) {
	bot, err := lg.lb.GetClient()
	if err != nil {
		return nil, err
	}
	res, err := bot.GetProfile(id).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the line profile; id = %s, err = %w", id, err)
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return entity.NewUser(res), nil
}

func (lg *LineGateway) GetQuotaComsuption() (int64, error) {
	bot, err := lg.lb.GetClient()
	if err != nil {
		return 0, err
	}
	res, err := bot.GetMessageQuotaConsumption().Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get the message quota consumption; err = %w", err)
	}
	conf.Log.Info("Successfully get the message quota consumption", zap.Any("res", res))
	return res.TotalUsage, nil
}

func (lg *LineGateway) GetFileContent(messageId string) (io.ReadCloser, error) {
	bot, err := lg.lb.GetClient()
	if err != nil {
		return nil, err
	}
	content, err := bot.GetMessageContent(messageId).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file content from LINE; messageId = %s, err = %w", messageId, err)
	}
	defer content.Content.Close()
	conf.Log.Info("Successfully download the file", zap.Int64("size (byte)", content.ContentLength))
	return content.Content, nil
}
