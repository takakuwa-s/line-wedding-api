package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load("./environments/dev.env")
	if err != nil {
		conf.Log.Error("Error loading .env file", zap.Any("err", err))
		return
	}

	flag.Parse()
	if flag.Arg(0) == "config" {
		conf.Log.Info("start configuration")
		if err := conf.GetRichmenuList(); err != nil {
			conf.Log.Error("Failed to get rich menu list", zap.Any("err", err))
			return
		}
		// if err := conf.DeleteRichmenu("richmenu-3bbe1606ee5c031bfc0313a7850874ea"); err != nil {
		// 	conf.Log.Error("Failed to delete the rich menu", zap.Any("err", err))
		// 	return
		// }
		if err := conf.CreateRichmenu(); err != nil {
			conf.Log.Error("Failed to create the rich menu", zap.Any("err", err))
			return
		}
		conf.Log.Info("complete configuration")
	} else {
		router := InitializeRouter()
		router.Init()
	}
}
