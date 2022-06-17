package main

import (
	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		conf.Log.Error("Error loading .env file", zap.Error(err))
		return
	}
	scheduler := InitializeScheduler()
	scheduler.Init()

	router := InitializeRouter()
	router.Init()
}
