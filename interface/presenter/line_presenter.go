package presenter

import (
	conf "github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"go.uber.org/zap"
)

type LinePresenter struct {
	bot *linebot.Client
}

// コンストラクタ
func NewLinePresenter(bot *linebot.Client) *LinePresenter {
	return &LinePresenter{bot: bot}
}

func (lp *LinePresenter) ReplyMessage(message dto.ReplyMessage) {
	sendingMessage := make([]linebot.SendingMessage, len(message.Messages))
	for i, m := range message.Messages {
		switch m["type"].(string) {
		case "text":
			sendingMessage[i] = linebot.NewTextMessage(m["text"].(string))
		case "sticker":
			sendingMessage[i] = linebot.NewStickerMessage(m["packageID"].(string), m["stickerID"].(string))
		}
		if m["quickReplyItems"] != nil {
			sendingMessage[i] = sendingMessage[i].WithQuickReplies(createQuickReplyItems(m))
		}
		if m["emojis"] != nil {
			for _, e := range createAllEmojis(m) {
				sendingMessage[i] = sendingMessage[i].AddEmoji(e)
			}
		}
		sendingMessage[i] = sendingMessage[i].WithSender(linebot.NewSender("くま", "https://drive.google.com/uc?export=view&id=1316g1ZQWfmffVy_uAVc_4UCxkG-9PTPq"))
	}
	if _, err := lp.bot.ReplyMessage(message.ReplyToken, sendingMessage...).Do(); err != nil {
		conf.Log.Error("Failed to send the reply message", zap.Any("err", err), zap.Any("messages", sendingMessage))
	}
}

func createQuickReplyItems(m map[string]interface{}) *linebot.QuickReplyItems {
	items := m["quickReplyItems"].([]interface{})
	buttons := make([]*linebot.QuickReplyButton, len(items))
	for i, item := range items {
		action := item.(map[string]interface{})["action"].(map[string]interface{})
		var lineAction linebot.QuickReplyAction
		switch action["type"] {
		case "message":
			lineAction = linebot.NewMessageAction(action["label"].(string), action["text"].(string))
		case "cameraRoll":
			lineAction = linebot.NewCameraRollAction(action["label"].(string))
		}
		imageURL := item.(map[string]interface{})["imageURL"]
		buttons[i] = linebot.NewQuickReplyButton(imageURL.(string), lineAction)
	}
	return linebot.NewQuickReplyItems(buttons...)
}

func createAllEmojis(m map[string]interface{}) []*linebot.Emoji {
	emojis := m["emojis"].([]interface{})
	res := make([]*linebot.Emoji, len(emojis))
	for i, emoji := range emojis {
		e := emoji.(map[string]interface{})
		res[i] = linebot.NewEmoji(int(e["index"].(float64)), e["productId"].(string), e["emojiId"].(string))
	}
	return res
}