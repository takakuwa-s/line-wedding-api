package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction()
)

func main() {
	err := godotenv.Load("./environments/dev.env")
	if err != nil {
		logger.Error("Error loading .env file", zap.Any("err", err))
	}
	router := InitializeRouter()
	router.Init()
}
