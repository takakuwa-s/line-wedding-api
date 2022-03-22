package presenter

import (
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/dto"

	"encoding/json"
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

func (lp *LinePresenter) ReplyMessage(message *dto.ReplyMessage) {
	sendingMessage := make([]linebot.SendingMessage, len(message.Messages))
	for i, m := range message.Messages {
		switch m["type"].(string) {
		case "text":
			sendingMessage[i] = linebot.NewTextMessage(m["text"].(string))
		case "sticker":
			sendingMessage[i] = linebot.NewStickerMessage(m["packageID"].(string), m["stickerID"].(string))
		case "image":
			sendingMessage[i] = linebot.NewImageMessage(m["originalContentUrl"].(string), m["previewImageUrl"].(string))
		case "video":
			sendingMessage[i] = linebot.NewVideoMessage(m["originalContentUrl"].(string), m["previewImageUrl"].(string))
		case "audio":
			sendingMessage[i] = linebot.NewAudioMessage(m["originalContentUrl"].(string), int(m["duration"].(float64)))
		case "location":
			sendingMessage[i] = linebot.NewLocationMessage(m["title"].(string), m["address"].(string), m["latitude"].(float64), m["longitude"].(float64))
		case "imagemap":
			sendingMessage[i] = createImageMapMessage(m)
		case "template":
			sendingMessage[i] = createTemplateMessage(m)
		case "flex":
			sendingMessage[i] = createFlexMessage(m)
		}
		if m["quickReply"] != nil {
			sendingMessage[i] = sendingMessage[i].WithQuickReplies(createQuickReplyItems(m))
		}
		if m["sender"] != nil {
			s := m["sender"].(map[string]interface{})
			sendingMessage[i] = sendingMessage[i].WithSender(linebot.NewSender(s["name"].(string), s["iconUrl"].(string)))
		}
		if m["emojis"] != nil {
			for _, e := range createAllEmojis(m) {
				sendingMessage[i] = sendingMessage[i].AddEmoji(e)
			}
		}
	}
	if _, err := lp.bot.ReplyMessage(message.ReplyToken, sendingMessage...).Do(); err != nil {
		conf.Log.Error("Failed to send the reply message", zap.Any("err", err), zap.Any("messages", sendingMessage))
	}
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

func createQuickReplyItems(m map[string]interface{}) *linebot.QuickReplyItems {
	items := m["quickReply"].(map[string]interface{})["items"].([]interface{})
	buttons := make([]*linebot.QuickReplyButton, len(items))
	for i, item := range items {
		action := createAction(item.(map[string]interface{})["action"].(map[string]interface{}))
		imageURL := item.(map[string]interface{})["imageURL"]
		buttons[i] = linebot.NewQuickReplyButton(imageURL.(string), action.(linebot.QuickReplyAction))
	}
	return linebot.NewQuickReplyItems(buttons...)
}

func createAction(m map[string]interface{}) linebot.Action {
	var action linebot.Action
	switch m["type"] {
	case "postback":
		action = linebot.NewPostbackAction(m["label"].(string), m["data"].(string), "", m["displayText"].(string))
	case "message":
		action = linebot.NewMessageAction(m["label"].(string), m["text"].(string))
	case "uri":
		action = linebot.NewURIAction(m["label"].(string), m["uri"].(string))
	case "datetimepicker":
		action = linebot.NewDatetimePickerAction(m["label"].(string), m["data"].(string),
		m["mode"].(string), m["initial"].(string), m["max"].(string), m["min"].(string))
	case "camera":
		action = linebot.NewCameraAction(m["label"].(string))
	case "cameraRoll":
		action = linebot.NewCameraRollAction(m["label"].(string))
	case "location":
		action = linebot.NewLocationAction(m["label"].(string))
	}
	return action
}

func createImageMapMessage(m map[string]interface{}) *linebot.ImagemapMessage {
	actions := make([]linebot.ImagemapAction, len(m["actions"].([]interface{})))
	for i, action := range m["actions"].([]interface{}) {
		a := action.(map[string]interface{})
		area := linebot.ImagemapArea{
			X: int(a["area"].(map[string]interface{})["x"].(float64)),
			Y: int(a["area"].(map[string]interface{})["y"].(float64)),
			Width: int(a["area"].(map[string]interface{})["width"].(float64)),
			Height: int(a["area"].(map[string]interface{})["height"].(float64)),
		}
		switch a["type"] {
		case "uri":
			actions[i] = linebot.NewURIImagemapAction(a["label"].(string), a["linkUri"].(string), area)
		case "message":
			actions[i] = linebot.NewMessageImagemapAction(a["label"].(string), a["text"].(string), area)
		}
	}
	baseSize := linebot.ImagemapBaseSize{
		Width: int(m["baseSize"].(map[string]interface{})["width"].(float64)),
		Height: int(m["baseSize"].(map[string]interface{})["height"].(float64)),
	}
	return linebot.NewImagemapMessage(m["baseUrl"].(string), m["altText"].(string), baseSize, actions...)
}

func createTemplateMessage(m map[string]interface{}) *linebot.TemplateMessage {
	t := m["template"].(map[string]interface{})
	var template linebot.Template
	switch t["type"].(string) {
	case "buttons":
		actions := createTemplateActions(t["actions"].([]interface{}))
		template = linebot.NewButtonsTemplate(t["thumbnailImageUrl"].(string), t["title"].(string), t["text"].(string), actions...)
	case "confirm":
		actions := createTemplateActions(t["actions"].([]interface{}))
		template = linebot.NewConfirmTemplate(t["text"].(string), actions[0], actions[1])
	case "carousel":
		columns := make([]*linebot.CarouselColumn, len(t["columns"].([]interface{})))
		for i, col := range t["columns"].([]interface{}) {
			c := col.(map[string]interface{})
			actions := createTemplateActions(c["actions"].([]interface{}))
			columns[i] = linebot.NewCarouselColumn(c["thumbnailImageUrl"].(string), c["title"].(string), c["text"].(string), actions...)
		}
		template = linebot.NewCarouselTemplate(columns...)
	case "image_carousel":
		columns := make([]*linebot.ImageCarouselColumn, len(t["columns"].([]interface{})))
		for i, col := range t["columns"].([]interface{}) {
			c := col.(map[string]interface{})
			action := createAction(c["action"].(map[string]interface{})).(linebot.TemplateAction)
			columns[i] = linebot.NewImageCarouselColumn(c["imageUrl"].(string), action)
		}
		template = linebot.NewImageCarouselTemplate(columns...)
	}
	return linebot.NewTemplateMessage(m["altText"].(string), template)
}

func createTemplateActions(m []interface{}) []linebot.TemplateAction {
	actions := make([]linebot.TemplateAction, len(m))
	for i, a := range m {
		actions[i] = createAction(a.(map[string]interface{})).(linebot.TemplateAction)
	}
	return actions
}

func createFlexMessage(m map[string]interface{}) *linebot.FlexMessage {
	c := m["contents"].(map[string]interface{})
	b, err :=	json.Marshal(c)
	if err != nil {
		conf.Log.Error("Failed to convert contents to byte", zap.Any("err", err))
	}
	flexContainer, err := linebot.UnmarshalFlexMessageJSON(b)
	if err != nil {
		conf.Log.Error("Failed to convert byte to flexContainer", zap.Any("err", err))
	}
	return linebot.NewFlexMessage(m["altText"].(string), flexContainer)
}