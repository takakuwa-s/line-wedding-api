package ipresenter

import (
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type IPresenter interface {
	MulticastMessage(message *dto.MulticastMessage, botType dto.BotType) error
	ReplyMessage(message *dto.ReplyMessage, botType dto.BotType) error
}