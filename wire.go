//+build wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"

	"github.com/google/wire"
)

func InitializeRouter() *driver.Router {
	wire.Build(

		// driver
		driver.NewRouter,

		// controller
		controller.NewLineController,

		// gateway
		gateway.NewMessageRepository,
		wire.Bind(new(igateway.IMessageRepository), new(*gateway.MessageRepository)),

		// presenter
		presenter.NewLinePresenter,
		wire.Bind(new(ipresenter.IPresenter), new(*presenter.LinePresenter)),

		// usecase
		usecase.NewMessageHandler,
	)
	return nil
}