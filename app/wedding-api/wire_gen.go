// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

// Injectors from wire.go:

func InitializeRouter() *driver.WeddingRouter {
	commonRouter := driver.NewCommonRouter()
	lineBot := dto.NewLineBot()
	messageRepository := gateway.NewMessageRepository()
	lineGateway := gateway.NewLineGateway(lineBot)
	firestore := dto.NewFirestore()
	commonRepository := gateway.NewCommonRepository(firestore)
	userRepository := gateway.NewUserRepository(commonRepository, firestore)
	fileRepository := gateway.NewFileRepository(commonRepository, firestore)
	imageSetRepository := gateway.NewImageSetRepository(commonRepository, firestore)
	fileUploadGateway := gateway.NewFileUploadGateway()
	linePresenter := presenter.NewLinePresenter(lineBot)
	linePushUsecase := usecase.NewLinePushUsecase(messageRepository, userRepository, linePresenter, lineGateway)
	lineReplyUsecase := usecase.NewLineReplyUsecase(messageRepository, lineGateway, userRepository, fileRepository, imageSetRepository, fileUploadGateway, linePushUsecase, linePresenter)
	lineBotController := controller.NewLineBotController(lineBot, lineReplyUsecase)
	binaryRepository := gateway.NewBinaryRepository(firestore)
	apiUsecase := usecase.NewApiUsecase(userRepository, lineGateway, fileRepository, binaryRepository)
	initApiController := controller.NewInitApiController(apiUsecase)
	userApiController := controller.NewUserApiController(apiUsecase)
	fileApiController := controller.NewFileApiController(apiUsecase)
	lineApiController := controller.NewLineApiController(linePushUsecase)
	weddingRouter := driver.NewWeddingRouter(commonRouter, lineBotController, initApiController, userApiController, fileApiController, lineApiController)
	return weddingRouter
}
