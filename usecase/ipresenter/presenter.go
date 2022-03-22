package ipresenter

import (
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
)

type IPresenter interface {
	ReplyMessage(*dto.ReplyMessage)
}