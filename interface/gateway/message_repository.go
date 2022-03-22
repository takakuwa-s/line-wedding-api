package gateway

import (
	"encoding/json"
	"io/ioutil"
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type MessageRepository struct {
	textReplyMessages map[string][]map[string]interface{}
	messages map[string][]map[string]interface{}
}

// Newコンストラクタ
func NewMessageRepository() *MessageRepository {
	textReplyMessages := readJson("./resource/text_reply_message.json")
	messages := readJson("./resource/message.json")
	return &MessageRepository{textReplyMessages: textReplyMessages, messages: messages}
}

func readJson(path string) map[string][]map[string]interface{} {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		conf.Log.Error("Failed to read the message.json", zap.Any("err", err), zap.String("path", path))
	}
	var obj map[string][]map[string]interface{}
	if err = json.Unmarshal(b, &obj); err != nil {
		conf.Log.Error("Failed to parses the JSON-encoded data", zap.Any("err", err))
	}
	return obj
}

func (mp *MessageRepository) FindReplyMessage(text string) []map[string]interface{} {
	ms := mp.textReplyMessages[text]
	if len(ms) > 0 {
		conf.Log.Info("Successfully find the textReplyMessages data", zap.Any("data", ms))
	}
	return ms
}

func (mp *MessageRepository) FindImageMessage() []map[string]interface{} {
	return mp.findMessage("image")
}

func (mp *MessageRepository) FindGroupMessage() []map[string]interface{} {
	return mp.findMessage("group")
}

func (mp *MessageRepository) FindFollowMessage(displayName string) []map[string]interface{} {
	ms := mp.findMessage("follow")
	ms[0]["text"] = fmt.Sprintf(ms[0]["text"].(string), displayName)
	return ms
}

func (mp *MessageRepository) findMessage(m string) []map[string]interface{} {
	ms := mp.messages[m]
	if len(ms) > 0 {
		conf.Log.Info("Successfully find the messages data", zap.Any("message", ms))
	}
	return ms
}

