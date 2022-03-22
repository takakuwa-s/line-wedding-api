package dto

type TextMessage struct {
	ReplyToken string
	Text string
}

func NewTextMessage(ReplyToken, Text string) *TextMessage {
	return &TextMessage{
		ReplyToken : ReplyToken,
		Text : Text,
	}
}

type FileMessage struct {
	ReplyToken string
	MessageID string
}

func NewFileMessage(ReplyToken, MessageID string) *FileMessage {
	return &FileMessage{
		ReplyToken : ReplyToken,
		MessageID: MessageID,
	}
}

type FollowMessage struct {
	ReplyToken string
	SenderUserID string
}

func NewFollowMessage(ReplyToken, SenderUserID string) *FollowMessage {
	return &FollowMessage{
		ReplyToken : ReplyToken,
		SenderUserID: SenderUserID,
	}
}

type GroupMessage struct {
	ReplyToken string
}

func NewGroupMessage(ReplyToken string) *GroupMessage {
	return &GroupMessage{
		ReplyToken : ReplyToken,
	}
}