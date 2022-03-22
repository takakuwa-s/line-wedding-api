package dto

type TextMessage struct {
	ReplyToken string
	Text string
}

func NewTextMessage(replyToken, text string) *TextMessage {
	return &TextMessage{
		ReplyToken : replyToken,
		Text : text,
	}
}

type FileMessage struct {
	ReplyToken string
	File *File
}

func NewFileMessage(replyToken, messageId string) *FileMessage {
	file := NewFile(messageId)
	return &FileMessage{
		ReplyToken : replyToken,
		File: file,
	}
}

type FollowMessage struct {
	ReplyToken string
	SenderUserId string
}

func NewFollowMessage(replyToken, senderUserId string) *FollowMessage {
	return &FollowMessage{
		ReplyToken : replyToken,
		SenderUserId: senderUserId,
	}
}

type PostbackMessage struct {
	ReplyToken string
	Data string
	Params interface{}
}

func NewPostbackMessage(replyToken, data string, params interface{}) *PostbackMessage {
	return &PostbackMessage{
		ReplyToken : replyToken,
		Data: data,
		Params: params,
	}
}

type GroupMessage struct {
	ReplyToken string
}

func NewGroupMessage(replyToken string) *GroupMessage {
	return &GroupMessage{
		ReplyToken : replyToken,
	}
}