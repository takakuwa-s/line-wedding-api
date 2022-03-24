package lineconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"go.uber.org/zap"
)

func GetRichmenuList(botType dto.BotType) error {
	bot, err := getBot(botType)
	if err != nil {
		return fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	res, err := bot.GetRichMenuList().Do()
	if err != nil {
		return fmt.Errorf("failed to get the richmenu list; err = %w", err)
	}
	conf.Log.Info("richmenu list size", zap.Int("len", len(res)))
	for _, richMenu := range res {
		conf.Log.Info("richMenu", zap.Any("richMenu", richMenu))
	}
	return nil
}

func DeleteRichmenu(richMenuID string, botType dto.BotType) error {
	bot, err := getBot(botType)
	if err != nil {
		return fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	if _, err := bot.DeleteRichMenu(richMenuID).Do(); err != nil {
		return fmt.Errorf("failed to delete the richmenu; err = %w", err)
	}
	conf.Log.Info("richmenu is deleted", zap.String("richMenuID", richMenuID))
	return nil
}

func CreateRichmenu(botType dto.BotType) error {
	bot, err := getBot(botType)
	if err != nil {
		return fmt.Errorf("failed to get the line bot client; err = %w", err)
	}
	path := fmt.Sprintf("./conf/resource/%s/richmenu.json", botType)
	b, err := ioutil.ReadFile(path)
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
	conf.Log.Info("rich menu call", zap.Any("CreateRichMenuCall", res))
	path = fmt.Sprintf("./conf/resource/%s/richmenu.png", botType)
	if _, err := bot.UploadRichMenuImage(res.RichMenuID, path).Do(); err != nil {
		return fmt.Errorf("failed to upload the image for the rich menu; err = %w", err)
	}
	if _, err := bot.SetDefaultRichMenu(res.RichMenuID).Do(); err != nil {
		return fmt.Errorf("failed to set the menu default; err = %w", err)
	}
	return nil
}

func getBot(botType dto.BotType) (*linebot.Client, error) {
	switch botType {
	case dto.WeddingBotType:
		return dto.NewWeddingLineBot().Client, nil
	case dto.AdminBotType:
		return dto.NewAdminLineBot().Client, nil
	default:
		return nil, fmt.Errorf("unknown bot type; %s", botType)
	}
}
