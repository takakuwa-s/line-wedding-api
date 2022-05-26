package main

import (
	"flag"

	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/conf/lineconf"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		conf.Log.Error("Error loading .env file", zap.Any("err", err))
		return
	}

	flag.Parse()
	if flag.Arg(0) == "config" {
		menuList, err := lineconf.GetRichmenuList();
		if err != nil {
			conf.Log.Error("Failed to get rich menu list", zap.Any("err", err))
			return
		}
		aliasList, err := lineconf.GetRichmenuAliasList()
		if err != nil {
			conf.Log.Error("Failed to get rich menu alias list", zap.Any("err", err))
			return
		}
		for _, menu := range menuList {
			if err := lineconf.DeleteRichmenu(menu.RichMenuID); err != nil {
				conf.Log.Error("Failed to delete the rich menu", zap.Any("err", err))
				return
			}
		}
		for _, alias := range aliasList {
			if err := lineconf.DeleteRichmenuAlias(alias.RichMenuAliasID); err != nil {
				conf.Log.Error("Failed to delete the rich menu", zap.Any("err", err))
				return
			}
		}
		if err := lineconf.CreateRichmenu(); err != nil {
			conf.Log.Error("Failed to create the rich menu", zap.Any("err", err))
			return
		}
		conf.Log.Info("complete configuration")
	} else {
		router := InitializeRouter()
		router.Init()
	}
}
