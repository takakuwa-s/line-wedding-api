package igateway

type IMessageRepository interface {
	FindReplyMessage(text string) []map[string]interface{}
	FindGroupMessage() []map[string]interface{}
	FindFollowMessage(displayName string) []map[string]interface{}
	FindImageMessage() []map[string]interface{}
}