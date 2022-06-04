package igateway

type IMessageRepository interface {
	FindReplyMessage(text string) []map[string]interface{}
	FindMessageByKey(key string) []map[string]interface{}
}
