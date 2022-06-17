//go:build wireinject
// +build wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/igateway"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"

	"github.com/google/wire"
)

func InitializeRouter() *driver.FileUploadRouter {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewFileUploadRouter,
		driver.NewCommonRouter,

		// controller
		controller.NewFileUploadController,

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
		usecase.NewFileUploadUsecase,
	)
	return nil
}

func InitializeScheduler() *driver.FileUploadScheduler {
	wire.Build(
		// dto
		dto.NewLineBot,
		dto.NewFirestore,

		// driver
		driver.NewFileUploadScheduler,

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
		usecase.NewFileUploadUsecase,
	)
	return nil
}
