package gateway

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/conf"
)

type LineRepository struct {
	bot *linebot.Client
}

// Newコンストラクタ
func NewLineRepository(bot *linebot.Client) *LineRepository {
	return &LineRepository{bot: bot}
}

func (lr *LineRepository) FindUserById(id string) *dto.User {
	res, err := lr.bot.GetProfile(id).Do();
	if err != nil {
		conf.Log.Error("Failed to get the line profile", zap.Any("err", err))
	}
	conf.Log.Info("Successfully find user", zap.Any("res", res))
	return &dto.User{
		UserId : res.UserID,
		DisplayName : res.DisplayName,
		Language : res.Language,
		PictureUrl : res.PictureURL,
		StatusMessage : res.StatusMessage,
	}
}