package gateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type MessageRepository struct {
	replyMessagePath string
	messagePath      string
}

// Newコンストラクタ
func NewMessageRepository() *MessageRepository {
	p1 := os.Getenv("REPLY_MESSAGE_FILE_PATH")
	p2 := os.Getenv("MESSAGE_FILE_PATH")
	return &MessageRepository{replyMessagePath: p1, messagePath: p2}
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
	ms := mp.readJson(mp.replyMessagePath)
	m := ms[text]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the reply messages data", zap.Any("data", m))
	}
	return m
}

func (mp *MessageRepository) FindMessageByKey(key string) []map[string]interface{} {
	ms := mp.readJson(mp.messagePath)
	m := ms[key]
	if len(m) > 0 {
		conf.Log.Info("Successfully find the messages data", zap.Any("message", m))
	}
	return m
}
