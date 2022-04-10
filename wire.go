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
		dto.NewWeddingLineBot,
		dto.NewAdminLineBot,
		dto.NewFirestore,

		// driver
		driver.NewRouter,

		// controller
		controller.NewWeddingLineController,
		controller.NewAdminLineController,

		// gateway
		gateway.NewMessageRepository,
		wire.Bind(new(igateway.IMessageRepository), new(*gateway.MessageRepository)),
		gateway.NewLineRepository,
		wire.Bind(new(igateway.ILineRepository), new(*gateway.LineRepository)),
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
		usecase.NewWeddingReplyUsecase,
		usecase.NewWeddingPushUsecase,
		usecase.NewAdminReplyUsecase,
		usecase.NewAdminPushUsecase,
		usecase.NewCommonUtils,
	)
	return nil
}