//+build wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"

	"github.com/google/wire"
)

func InitializeRouter() *driver.Router {
	wire.Build(

		// driver
		driver.NewRouter,

		// controller
		controller.NewLineController,

		// presenter
		presenter.NewLinePresenter,

		// usecase
		usecase.NewMessageHandler,
		wire.Bind(new(ipresenter.IPresenter), new(*presenter.LinePresenter)),
	)
	return nil
}