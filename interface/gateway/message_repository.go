package gateway

import (
	"encoding/json"
	"io/ioutil"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type MessageRepository struct {
	data map[string][]map[string]interface{}
}

// Newコンストラクタ
func NewMessageRepository() *MessageRepository {
	b, err := ioutil.ReadFile("./resource/message.json")
	if err != nil {
		conf.Log.Error("Failed to read the message.json", zap.Any("err", err))
	}
	var obj map[string][]map[string]interface{}
	if err = json.Unmarshal(b, &obj); err != nil {
		conf.Log.Error("Failed to parses the JSON-encoded data", zap.Any("err", err))
	}
	return &MessageRepository{data: obj}
}

func (mp *MessageRepository) FindReplyMessage(text string) []map[string]interface{} {
	ms := mp.data[text]
	if len(ms) > 0 {
		conf.Log.Info("Successfully find the data", zap.Any("data", ms))
	}
	return ms
}