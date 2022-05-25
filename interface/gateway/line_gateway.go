package gateway

import (
	"fmt"
	"io"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type LineGateway struct {
	wlb *dto.WeddingLineBot
	alb *dto.AdminLineBot
}

// Newコンストラクタ
func NewLineGateway(wlb *dto.WeddingLineBot, alb *dto.AdminLineBot) *LineGateway {
	return &LineGateway{wlb: wlb, alb: alb}
}

func (lg *LineGateway) GetUserProfileById(id string, botType dto.BotType) (*entity.User, error) {
	bot, err := lg.getBot(botType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	res, err := bot.GetProfile(id).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the line profile; err = %w", err)
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return entity.NewUser(res), nil
}

func (lg *LineGateway) GetQuotaComsuption(botType dto.BotType) (int64, error) {
	bot, err := lg.getBot(botType)
	if err != nil {
		return 0, fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	res, err := bot.GetMessageQuotaConsumption().Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get the message quota consumption; err = %w", err)
	}
	conf.Log.Info("Successfully get the message quota consumption", zap.Any("res", res))
	return res.TotalUsage, nil
}

func (lg *LineGateway) GetFileContent(botType dto.BotType, messageId string) (io.ReadCloser, error) {
	bot, err := lg.getBot(botType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	content, err := bot.GetMessageContent(messageId).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the file content from LINE; err = %w", err)
	}
	defer content.Content.Close()
	conf.Log.Info("Successfully download the file", zap.String("ContentType", content.ContentType), zap.Int64("size (byte)", content.ContentLength))
	return content.Content, nil
}

func (lg *LineGateway) getBot(botType dto.BotType) (*linebot.Client, error) {
	switch botType {
	case dto.WeddingBotType:
		return lg.wlb.Client, nil
	case dto.AdminBotType:
		return lg.alb.Client, nil
	default:
		return nil, fmt.Errorf("unknown bot type; %s", botType)
	}
}
