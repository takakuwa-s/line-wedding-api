package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/conf/lineconf"
	"github.com/takakuwa-s/line-wedding-api/dto"
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
		botType := dto.AdminBotType
		conf.Log.Info("start configuration", zap.Any("botType", botType))
		if err := lineconf.GetRichmenuList(botType); err != nil {
			conf.Log.Error("Failed to get rich menu list", zap.Any("err", err))
			return
		}
		// if err := lineconf.DeleteRichmenu("richmenu-617f4a4dc1494564ecae3113e43c2285", botType); err != nil {
		// 	conf.Log.Error("Failed to delete the rich menu", zap.Any("err", err))
		// 	return
		// }
		// if err := lineconf.CreateRichmenu(botType); err != nil {
		// 	conf.Log.Error("Failed to create the rich menu", zap.Any("err", err))
		// 	return
		// }
		conf.Log.Info("complete configuration")
	} else {
		router := InitializeRouter()
		router.Init()
	}
}
