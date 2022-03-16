package igateway

type IMessageRepository interface {
	FindReplyMessage(string) string
}