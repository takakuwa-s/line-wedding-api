package gateway

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type LineRepository struct {
	bot *linebot.Client
}

// Newコンストラクタ
func NewLineRepository(bot *linebot.Client) *LineRepository {
	return &LineRepository{bot: bot}
}

func (lr *LineRepository) FindUserById(id string) (*entity.User, error) {
	res, err := lr.bot.GetProfile(id).Do();
	if err != nil {
		return nil, fmt.Errorf("failed to get the line profile; err = %w", err)
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return entity.NewUser(res), nil
}

func (lr *LineRepository) GetQuotaComsuption() (int, error) {
	res, err := lr.bot.GetMessageQuotaConsumption().Do()
	if err != nil {
		return 0, fmt.Errorf("failed to get the message quota consumption; err = %w", err)
	}
	conf.Log.Info("Successfully get the message quota consumption", zap.Any("res", res))
	return 1, nil
}