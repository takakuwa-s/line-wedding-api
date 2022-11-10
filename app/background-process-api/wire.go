//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

func InitializeRouter() *driver.BackgroundProcessRouter {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewBackgroundProcessRouter,
		driver.NewCommonRouter,

		// controller
		controller.NewBackgroundProcessController,
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
		gateway.NewFaceGateway,
		wire.Bind(new(igateway.IFaceGateway), new(*gateway.FaceGateway)),
		gateway.NewLineGateway,
		wire.Bind(new(igateway.ILineGateway), new(*gateway.LineGateway)),
		gateway.NewSlideShowGateway,
		wire.Bind(new(igateway.ISlideShowGateway), new(*gateway.SlideShowGateway)),
		gateway.NewFSlideShowRepository,
		wire.Bind(new(igateway.ISlideShowRepository), new(*gateway.SlideShowRepository)),

		// presenter
		presenter.NewLinePresenter,
		wire.Bind(new(ipresenter.IPresenter), new(*presenter.LinePresenter)),

		// usecase
		usecase.NewBackgroundProcessUsecase,
		usecase.NewSlideShowUsecase,
		usecase.NewLinePushUsecase,
	)
	return nil
}

func InitializeScheduler() *driver.BackgroundProcessScheduler {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewBackgroundProcessScheduler,

		// gateway
		gateway.NewCommonRepository,
		gateway.NewFileRepository,
		wire.Bind(new(igateway.IFileRepository), new(*gateway.FileRepository)),
		gateway.NewBinaryRepository,
		wire.Bind(new(igateway.IBinaryRepository), new(*gateway.BinaryRepository)),
		gateway.NewFaceGateway,
		wire.Bind(new(igateway.IFaceGateway), new(*gateway.FaceGateway)),
		gateway.NewLineGateway,
		wire.Bind(new(igateway.ILineGateway), new(*gateway.LineGateway)),

		// usecase
		usecase.NewBackgroundProcessUsecase,
	)
	return nil
}
