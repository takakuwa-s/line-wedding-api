package ipresenter

import (
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type IPresenter interface {
	MulticastMessage(message *dto.MulticastMessage) error
	ReplyMessage(message *dto.ReplyMessage) error
}
