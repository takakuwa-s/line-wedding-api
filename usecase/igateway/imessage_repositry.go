package igateway

import "github.com/takakuwa-s/line-wedding-api/dto"

type IMessageRepository interface {
	FindReplyMessage(botType dto.BotType, text string) []map[string]interface{}
	FindMessageByKey(botType dto.BotType, key string) []map[string]interface{}
}