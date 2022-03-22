package dto

type ReplyMessage struct {
	ReplyToken string
	Messages []map[string]interface{}
}

func NewReplyMessage(ReplyToken string, Messages []map[string]interface{}) *ReplyMessage {
	return &ReplyMessage{
		ReplyToken : ReplyToken,
		Messages : Messages,
	}
} 