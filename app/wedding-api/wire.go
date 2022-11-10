//go:build wireinject
// +build wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"

	"github.com/google/wire"
)

func InitializeRouter() *driver.WeddingRouter {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewWeddingRouter,
		driver.NewCommonRouter,

		// controller
		controller.NewLineBotController,
		controller.NewInitApiController,
		controller.NewUserApiController,
		controller.NewFileApiController,
		controller.NewLineApiController,
		controller.NewSlideShowApiController,

		// gateway
		gateway.NewCommonRepository,
		gateway.NewMessageRepository,
		wire.Bind(new(igateway.IMessageRepository), new(*gateway.MessageRepository)),
		gateway.NewUserRepository,
		wire.Bind(new(igateway.IUserRepository), new(*gateway.UserRepository)),
		gateway.NewFileRepository,
		wire.Bind(new(igateway.IFileRepository), new(*gateway.FileRepository)),
		gateway.NewBinaryRepository,
		wire.Bind(new(igateway.IBinaryRepository), new(*gateway.BinaryRepository)),
		gateway.NewImageSetRepository,
		wire.Bind(new(igateway.IImageSetRepository), new(*gateway.ImageSetRepository)),
		gateway.NewLineGateway,
		wire.Bind(new(igateway.ILineGateway), new(*gateway.LineGateway)),
		gateway.NewBackgroundProcessGateway,
		wire.Bind(new(igateway.IBackgroundProcessGateway), new(*gateway.BackgroundProcessGateway)),
		gateway.NewSlideShowGateway,
		wire.Bind(new(igateway.ISlideShowGateway), new(*gateway.SlideShowGateway)),
		gateway.NewFSlideShowRepository,
		wire.Bind(new(igateway.ISlideShowRepository), new(*gateway.SlideShowRepository)),

		// presenter
		presenter.NewLinePresenter,
		wire.Bind(new(ipresenter.IPresenter), new(*presenter.LinePresenter)),

		// usecase
		usecase.NewLineReplyUsecase,
		usecase.NewLinePushUsecase,
		usecase.NewApiUsecase,
		usecase.NewSlideShowUsecase,
	)
	return nil
}
