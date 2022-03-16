package gateway

import (
	"fmt"
	"go.uber.org/zap"
	"encoding/json"
	"io/ioutil"
)

var (
	logger, _ = zap.NewProduction()
)

type MessageRepository struct {
	data map[string]interface{}
}

// Newコンストラクタ
func NewMessageRepository() *MessageRepository {
	b, err := ioutil.ReadFile("./resources/message.json")
	if err != nil {
		logger.Error("Failed to read the file", zap.Any("err", err))
	}
	var obj map[string]interface{}
	_ = json.Unmarshal(b, &obj)
	return &MessageRepository{data: obj}
}

func (mp *MessageRepository) FindReplyMessage(text string) string {
	message := mp.data[text]
	if fmt.Sprintf("%T", message) == "string" {
		logger.Info("Message is successfully found", zap.Any("message", message))
		return message.(string)
	}
	return ""
}