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
	backgroundProcessGateway := gateway.NewBackgroundProcessGateway(lineBot)
	linePresenter := presenter.NewLinePresenter(lineBot)
	linePushUsecase := usecase.NewLinePushUsecase(messageRepository, userRepository, linePresenter, lineGateway)
	slideShowGateway := gateway.NewSlideShowGateway()
	slideShowRepository := gateway.NewFSlideShowRepository(commonRepository, firestore)
	binaryRepository := gateway.NewBinaryRepository(firestore)
	slideShowUsecase := usecase.NewSlideShowUsecase(messageRepository, fileRepository, slideShowGateway, slideShowRepository, binaryRepository, linePushUsecase)
	lineReplyUsecase := usecase.NewLineReplyUsecase(messageRepository, lineGateway, userRepository, fileRepository, imageSetRepository, backgroundProcessGateway, linePushUsecase, slideShowUsecase, linePresenter)
	lineBotController := controller.NewLineBotController(lineBot, lineReplyUsecase)
	apiUsecase := usecase.NewApiUsecase(messageRepository, userRepository, lineGateway, fileRepository, binaryRepository, backgroundProcessGateway, linePushUsecase, slideShowUsecase)
	userApiController := controller.NewUserApiController(apiUsecase)
	fileApiController := controller.NewFileApiController(apiUsecase)
	lineApiController := controller.NewLineApiController(apiUsecase)
	slideShowApiController := controller.NewSlideShowApiController(slideShowUsecase)
	configApiController := controller.NewConfigApiController()
	weddingRouter := driver.NewWeddingRouter(commonRouter, lineBotController, userApiController, fileApiController, lineApiController, slideShowApiController, configApiController)
	return weddingRouter
}
