package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"go.uber.org/zap"
)

func GetRichmenuList() ([]*linebot.RichMenuResponse, error) {
	bot := dto.NewLineBot().Client
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the richmenu list; err = %w", err)
	}
	conf.Log.Info("richmenu list size", zap.Int("len", len(res)))
	for _, richMenu := range res {
		conf.Log.Info("richMenu", zap.Any("richMenu", richMenu))
	}
	return res, nil
}

func GetRichmenuAliasList() ([]*linebot.RichMenuAliasResponse, error) {
	bot := dto.NewLineBot().Client
	res, err := bot.GetRichMenuAliasList().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the richmenu alias list; err = %w", err)
	}
	conf.Log.Info("richmenu alias list size", zap.Int("len", len(res)))
	for _, alias := range res {
		conf.Log.Info("richMenu alias", zap.Any("richMenu", alias))
	}
	return res, nil
}

func DeleteRichmenu(richMenuID string) error {
	bot := dto.NewLineBot().Client
	if _, err := bot.DeleteRichMenu(richMenuID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu; err = %w", err)
	}
	conf.Log.Info("richmenu is deleted", zap.String("richMenuID", richMenuID))
	return nil
}

func DeleteRichmenuAlias(richMenuAliasID string) error {
	bot := dto.NewLineBot().Client
	if _, err := bot.DeleteRichMenuAlias(richMenuAliasID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu alias; err = %w", err)
	}
	conf.Log.Info("richmenu alias is deleted", zap.String("richMenuAliasID", richMenuAliasID))
	return nil
}

func createRichmenu(jsonPath, imagePath string, bot *linebot.Client) (string, error) {
	b, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to read the richmenu.json; err = %w", err)
	}
	var menu linebot.RichMenu
	if err = json.Unmarshal(b, &menu); err != nil {
		return "", fmt.Errorf("failed to parses the JSON-encoded data; err = %w", err)
	}
	res, err := bot.CreateRichMenu(menu).Do()
	if err != nil {
		return "", fmt.Errorf("failed to request to create the rich menu; err = %w", err)
	}
	conf.Log.Info("rich menu is successfully created", zap.Any("CreateRichMenuCall", res))
	if _, err := bot.UploadRichMenuImage(res.RichMenuID, imagePath).Do(); err != nil {
		return "", fmt.Errorf("failed to upload the image for the rich menu; err = %w", err)
	}
	return res.RichMenuID, nil
}

func CreateRichmenu() error {
	bot := dto.NewLineBot().Client
	bot.GetRichMenuAliasList()
	id1, err := createRichmenu("./resource/richmenu-1.json", "./resource/richmenu-1.png", bot)
	if err != nil {
		return err
	}
	id2, err := createRichmenu("./resource/richmenu-2.json", "./resource/richmenu-2.png", bot)
	if err != nil {
		return err
	}
	if _, err := bot.SetDefaultRichMenu(id1).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 1 as default; err = %w", err)
	}
	if _, err := bot.CreateRichMenuAlias("richmenu-alias-1", id1).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 1 as alias 1 ; err = %w", err)
	}
	if _, err := bot.CreateRichMenuAlias("richmenu-alias-2", id2).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 2 as alias 2 ; err = %w", err)
	}
	return nil
}
