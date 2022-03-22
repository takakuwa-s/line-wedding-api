package igateway

type IMessageRepository interface {
	FindReplyMessage(string) []map[string]interface{}
	FindGroupMessage() []map[string]interface{}
	FindFollowMessage() []map[string]interface{}
	FindImageMessage() []map[string]interface{}
}