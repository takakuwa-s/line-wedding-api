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

	lrc := NewLineRichmenuConf()
	menuList, err := lrc.GetRichmenuList()
	if err != nil {
		conf.Log.Error("Failed to get rich menu list", zap.Error(err))
		return
	}
	aliasList, err := lrc.GetRichmenuAliasList()
	if err != nil {
		conf.Log.Error("Failed to get rich menu alias list", zap.Error(err))
		return
	}
	for _, menu := range menuList {
		if err := lrc.DeleteRichmenu(menu.RichMenuID); err != nil {
			conf.Log.Error("Failed to delete the rich menu", zap.Error(err))
			return
		}
	}
	for _, alias := range aliasList {
		if err := lrc.DeleteRichmenuAlias(alias.RichMenuAliasID); err != nil {
			conf.Log.Error("Failed to delete the rich menu", zap.Error(err))
			return
		}
	}
	if err := lrc.CreateRichmenu(); err != nil {
		conf.Log.Error("Failed to create the rich menu", zap.Error(err))
		return
	}
	conf.Log.Info("complete configuration")
}
