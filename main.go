package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"

	"github.com/takakuwa-s/line-wedding-api/usecase/dto"
	"github.com/takakuwa-s/line-wedding-api/interface/gateway"
	"github.com/takakuwa-s/line-wedding-api/interface/presenter"
)

func main() {
	err := godotenv.Load("./environments/dev.env")
	if err != nil {
		conf.Log.Error("Error loading .env file", zap.Any("err", err))
	}

	flag.Parse()
	if flag.Arg(0) == "config" {
		// conf.Log.Info("start configuration")
		// conf.GetRichmenuList()
		// conf.DeleteRichmenu("richmenu-0892e23935cd215dc79c319c01b98383")
		// conf.CreateRichmenu()
		// conf.Log.Info("complete configuration")
		m := gateway.NewMessageRepository().FindReplyMessage("結婚式")
		bot := conf.NewLineBot()
		rm := dto.NewReplyMessage("", m)
		presenter.NewLinePresenter(bot).ReplyMessage(rm)
	} else {
		router := InitializeRouter()
		router.Init()
	}
}
