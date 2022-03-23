package conf

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"encoding/json"
	"io/ioutil"
	"go.uber.org/zap"
)

func GetRichmenuList() error {
	bot := NewLineBot()
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return fmt.Errorf("failed to get the richmenu list; err = %w", err)
	}
	Log.Info("richmenu list size", zap.Int("len", len(res)))
	for _, richMenu := range res {
		Log.Info("richMenu", zap.Any("richMenu", richMenu))
	}
	return nil
}

func DeleteRichmenu(richMenuID string) error {
	bot := NewLineBot()
	if _, err := bot.DeleteRichMenu(richMenuID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu; err = %w", err)
	}
	Log.Info("richmenu is deleted", zap.String("richMenuID", richMenuID))
	return nil
}

func CreateRichmenu() error {
	bot := NewLineBot()
	b, err := ioutil.ReadFile("./conf/resource/richmenu.json")
	if err != nil {
		return fmt.Errorf("failed to read the richmenu.json; err = %w", err)
	}
	var menu linebot.RichMenu
	if err = json.Unmarshal(b, &menu); err != nil {
		return fmt.Errorf("failed to parses the JSON-encoded data; err = %w", err)
	}
	res, err := bot.CreateRichMenu(menu).Do()
	if err != nil {
		return fmt.Errorf("failed to request to create the rich menu; err = %w", err)
	}
	Log.Info("rich menu call", zap.Any("CreateRichMenuCall", res))
	if _, err := bot.UploadRichMenuImage(res.RichMenuID, "./conf/resource/richmenu.png").Do(); err != nil {
		return fmt.Errorf("failed to upload the image for the rich menu; err = %w", err)
	}
	if _, err := bot.SetDefaultRichMenu(res.RichMenuID).Do(); err != nil {
		return fmt.Errorf("failed to set the menu default; err = %w", err)
	}
	return nil
}
