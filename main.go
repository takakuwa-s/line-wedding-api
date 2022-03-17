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
	}

	flag.Parse()
	if flag.Arg(0) == "config" {
		conf.Log.Info("start configuration")
		conf.GetRichmenuList()
		// conf.DeleteRichmenu("richmenu-c366d8ae6087d26dc448c4bca360b673")
		// conf.CreateRichmenu()
		conf.Log.Info("complete configuration")
	} else {
		router := InitializeRouter()
		router.Init()
	}
}
