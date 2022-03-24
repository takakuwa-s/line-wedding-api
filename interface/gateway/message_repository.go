package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
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

func (mp *MessageRepository) FindReplyMessage(botType dto.BotType, text string) []map[string]interface{} {
	path := fmt.Sprintf("./resource/%s/reply_message.json", botType)
	ms := mp.readJson(path)
	m := ms[text]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the reply messages data", zap.Any("data", m))
	}
	return m
}

func (mp *MessageRepository) FindMessageByKey(botType dto.BotType, key string) []map[string]interface{} {
	path := fmt.Sprintf("./resource/%s/message.json", botType)
	ms := mp.readJson(path)
	m := ms[key]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the messages data", zap.Any("message", ms))
	}
	return m
}

