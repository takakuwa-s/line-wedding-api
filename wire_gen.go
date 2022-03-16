// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/takakuwa-s/line-wedding-api/driver"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
	"github.com/takakuwa-s/line-wedding-api/usecase"
)

// Injectors from wire.go:

func InitializeRouter() *driver.Router {
	linePresenter := presenter.NewLinePresenter()
	messageRepository := gateway.NewMessageRepository()
	messageHandler := usecase.NewMessageHandler(linePresenter, messageRepository)
	lineController := controller.NewLineController(messageHandler)
	router := driver.NewRouter(lineController)
	return router
}
