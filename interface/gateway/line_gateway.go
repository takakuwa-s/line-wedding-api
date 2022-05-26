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
	res, err := lg.lb.Client.GetProfile(id).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the line profile; err = %w", err)
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return entity.NewUser(res), nil
}

func (lg *LineGateway) GetQuotaComsuption() (int64, error) {
	res, err := lg.lb.Client.GetMessageQuotaConsumption().Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get the message quota consumption; err = %w", err)
	}
	conf.Log.Info("Successfully get the message quota consumption", zap.Any("res", res))
	return res.TotalUsage, nil
}

func (lg *LineGateway) GetFileContent(messageId string) (io.ReadCloser, error) {
	content, err := lg.lb.Client.GetMessageContent(messageId).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file content from LINE; err = %w", err)
	}
	defer content.Content.Close()
	conf.Log.Info("Successfully download the file", zap.String("ContentType", content.ContentType), zap.Int64("size (byte)", content.ContentLength))
	return content.Content, nil
}
