package dto

type ReplyMessage struct {
	ReplyToken string
	Messages   []map[string]interface{}
}

func NewReplyMessage(replyToken string, messages []map[string]interface{}) *ReplyMessage {
	return &ReplyMessage{
		ReplyToken: replyToken,
		Messages:   messages,
	}
}
