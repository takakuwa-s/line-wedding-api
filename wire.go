//+build wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"github.com/takakuwa-s/line-wedding-api/usecase/ipresenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"

	"github.com/google/wire"
)

func InitializeRouter() *driver.Router {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewRouter,

		// controller
		controller.NewLineBotController,
		controller.NewInitApiController,
		controller.NewUserApiController,
		controller.NewFileApiController,
		controller.NewLineApiController,

		// gateway
		gateway.NewMessageRepository,
		wire.Bind(new(igateway.IMessageRepository), new(*gateway.MessageRepository)),
		gateway.NewFaceGateway,
		wire.Bind(new(igateway.IFaceGateway), new(*gateway.FaceGateway)),
		gateway.NewLineGateway,
		wire.Bind(new(igateway.ILineGateway), new(*gateway.LineGateway)),
		gateway.NewUserRepository,
		wire.Bind(new(igateway.IUserRepository), new(*gateway.UserRepository)),
		gateway.NewFileRepository,
		wire.Bind(new(igateway.IFileRepository), new(*gateway.FileRepository)),
		gateway.NewBinaryRepository,
		wire.Bind(new(igateway.IBinaryRepository), new(*gateway.BinaryRepository)),

		// presenter
		presenter.NewLinePresenter,
		wire.Bind(new(ipresenter.IPresenter), new(*presenter.LinePresenter)),

		// usecase
		usecase.NewLineReplyUsecase,
		usecase.NewLinePushUsecase,
		usecase.NewApiUsecase,
	)
	return nil
}