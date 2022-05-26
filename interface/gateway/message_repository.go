package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type MessageRepository struct {
}

// Newコンストラクタ
func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

func (mp *MessageRepository) readJson(path string) map[string][]map[string]interface{} {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to read the message.json; path = %s, err = %v", path, err))
	}
	var obj map[string][]map[string]interface{}
	if err = json.Unmarshal(b, &obj); err != nil {
		panic(fmt.Sprintf("failed to parses the JSON-encoded data; path = %s, err = %v", path, err))
	}
	return obj
}

func (mp *MessageRepository) FindReplyMessage(text string) []map[string]interface{} {
	ms := mp.readJson("./resource/reply_message.json")
	m := ms[text]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the reply messages data", zap.Any("data", m))
	}
	return m
}

func (mp *MessageRepository) FindMessageByKey(key string) []map[string]interface{} {
	ms := mp.readJson("./resource/message.json")
	m := ms[key]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the messages data", zap.Any("message", m))
	}
	return m
}

