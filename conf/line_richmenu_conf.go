package conf

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"

	"encoding/json"
	"io/ioutil"

	"go.uber.org/zap"
)

func GetRichmenuList() {
	bot := NewLineBot()
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		Log.Error("Failed to get the richmenu list", zap.Any("err", err))
	}
	Log.Info("richmenu list size", zap.Int("len", len(res)))
	for _, richMenu := range res {
		Log.Info("richMenu", zap.Any("richMenu", richMenu))
	}
}

func DeleteRichmenu(richMenuID string) {
	bot := NewLineBot()
	if _, err := bot.DeleteRichMenu(richMenuID).Do(); err != nil {
		Log.Error("Failed to delete the richmenu", zap.Any("err", err))
	}
	Log.Info("richmenu is deleted", zap.String("richMenuID", richMenuID))
}

func CreateRichmenu() {
	bot := NewLineBot()
	b, err := ioutil.ReadFile("./conf/resource/richmenu.json")
	if err != nil {
		Log.Error("Failed to read the richmenu.json", zap.Any("err", err))
	}
	var menu linebot.RichMenu
	if err = json.Unmarshal(b, &menu); err != nil {
		Log.Error("Failed to parses the JSON-encoded data", zap.Any("err", err))
	}
	res, err := bot.CreateRichMenu(menu).Do()
	if err != nil {
		Log.Error("Failed to request to create the rich menu", zap.Any("err", err))
	}
	Log.Info("rich menu call", zap.Any("CreateRichMenuCall", res))
	if _, err := bot.UploadRichMenuImage(res.RichMenuID, "./conf/resource/richmenu.png").Do(); err != nil {
		Log.Error("Failed to upload the image for the rich menu", zap.Any("err", err))
	}
	if _, err := bot.SetDefaultRichMenu(res.RichMenuID).Do(); err != nil {
		Log.Error("Failed to set the menu default", zap.Any("err", err))
	}
}
