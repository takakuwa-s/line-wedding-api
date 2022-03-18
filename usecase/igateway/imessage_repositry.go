package igateway

type IMessageRepository interface {
	FindReplyMessage(string) []map[string]interface{}
}