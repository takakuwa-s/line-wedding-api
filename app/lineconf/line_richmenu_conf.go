package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"go.uber.org/zap"
)

type LineRichmenuConf struct {
	client *linebot.Client
}

// コンストラクタ
func NewLineRichmenuConf() *LineRichmenuConf {
	client, err := dto.NewLineBot().GetClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to create the line bot instance; err = %v", err))
	}
	return &LineRichmenuConf{client: client}
}

func (lrc *LineRichmenuConf) GetRichmenuList() ([]*linebot.RichMenuResponse, error) {
	res, err := lrc.client.GetRichMenuList().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the richmenu list; err = %w", err)
	}
	conf.Log.Info("richmenu list size", zap.Int("len", len(res)))
	for _, richMenu := range res {
		conf.Log.Info("richMenu", zap.Any("richMenu", richMenu))
	}
	return res, nil
}

func (lrc *LineRichmenuConf) GetRichmenuAliasList() ([]*linebot.RichMenuAliasResponse, error) {
	res, err := lrc.client.GetRichMenuAliasList().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get the richmenu alias list; err = %w", err)
	}
	conf.Log.Info("richmenu alias list size", zap.Int("len", len(res)))
	for _, alias := range res {
		conf.Log.Info("richMenu alias", zap.Any("richMenu", alias))
	}
	return res, nil
}

func (lrc *LineRichmenuConf) DeleteRichmenu(richMenuID string) error {
	if _, err := lrc.client.DeleteRichMenu(richMenuID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu; err = %w", err)
	}
	conf.Log.Info("richmenu is deleted", zap.String("richMenuID", richMenuID))
	return nil
}

func (lrc *LineRichmenuConf) DeleteRichmenuAlias(richMenuAliasID string) error {
	if _, err := lrc.client.DeleteRichMenuAlias(richMenuAliasID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu alias; err = %w", err)
	}
	conf.Log.Info("richmenu alias is deleted", zap.String("richMenuAliasID", richMenuAliasID))
	return nil
}

func (lrc *LineRichmenuConf) createRichmenu(jsonPath, imagePath string) (string, error) {
	b, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return "", fmt.Errorf("failed to read the richmenu.json; err = %w", err)
	}
	var menu linebot.RichMenu
	if err = json.Unmarshal(b, &menu); err != nil {
		return "", fmt.Errorf("failed to parses the JSON-encoded data; err = %w", err)
	}
	liffUrl := os.Getenv("LIFF_URL")
	for i, area := range menu.Areas {
		if strings.Contains(area.Action.URI, "%s") {
			menu.Areas[i].Action.URI = fmt.Sprintf(area.Action.URI, liffUrl)
		}
	}
	res, err := lrc.client.CreateRichMenu(menu).Do()
	if err != nil {
		return "", fmt.Errorf("failed to request to create the rich menu; err = %w", err)
	}
	conf.Log.Info("rich menu is successfully created", zap.Any("CreateRichMenuCall", res))
	if _, err := lrc.client.UploadRichMenuImage(res.RichMenuID, imagePath).Do(); err != nil {
		return "", fmt.Errorf("failed to upload the image for the rich menu; err = %w", err)
	}
	return res.RichMenuID, nil
}

func (lrc *LineRichmenuConf) CreateRichmenu() error {
	id1, err := lrc.createRichmenu("./resource/richmenu-1.json", "./resource/richmenu-1.png")
	if err != nil {
		return err
	}
	id2, err := lrc.createRichmenu("./resource/richmenu-2.json", "./resource/richmenu-2.png")
	if err != nil {
		return err
	}
	if _, err := lrc.client.SetDefaultRichMenu(id1).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 1 as default; err = %w", err)
	}
	if _, err := lrc.client.CreateRichMenuAlias("richmenu-alias-1", id1).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 1 as alias 1 ; err = %w", err)
	}
	if _, err := lrc.client.CreateRichMenuAlias("richmenu-alias-2", id2).Do(); err != nil {
		return fmt.Errorf("failed to set the menu 2 as alias 2 ; err = %w", err)
	}
	return nil
}
