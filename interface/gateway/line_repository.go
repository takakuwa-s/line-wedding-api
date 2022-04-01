package gateway

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type LineRepository struct {
	wlb *dto.WeddingLineBot
	alb *dto.AdminLineBot
}

// Newコンストラクタ
func NewLineRepository(wlb *dto.WeddingLineBot, alb *dto.AdminLineBot) *LineRepository {
	return &LineRepository{wlb:wlb, alb:alb}
}

func (lr *LineRepository) GetUserProfileById(id string, botType dto.BotType) (*entity.User, error) {
	bot, err := lr.getBot(botType)
	if err != nil {
		return nil, fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	res, err := bot.GetProfile(id).Do();
	if err != nil {
		return nil, fmt.Errorf("failed to get the line profile; err = %w", err)
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return entity.NewUser(res), nil
}

func (lr *LineRepository) GetQuotaComsuption(botType dto.BotType) (int64, error) {
	bot, err := lr.getBot(botType)
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

func (lr *LineRepository) getBot(botType dto.BotType) (*linebot.Client, error) {
	switch botType {
	case dto.WeddingBotType:
		return lr.wlb.Client, nil
	case dto.AdminBotType:
		return lr.alb.Client, nil
	default:
		return nil, fmt.Errorf("unknown bot type; %s", botType)
	}
}