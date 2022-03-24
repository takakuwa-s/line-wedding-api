package dto

import (
	"time"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

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
	File *entity.File
}

func NewFileMessage(replyToken string, file *entity.File) *FileMessage {
	return &FileMessage{
		ReplyToken : replyToken,
		File: file,
	}
}

type FollowMessage struct {
	ReplyToken string
	SenderUserId string
	EventTime time.Time
}

func NewFollowMessage(replyToken, senderUserId string, eventTime time.Time) *FollowMessage {
	return &FollowMessage{
		ReplyToken : replyToken,
		SenderUserId: senderUserId,
		EventTime: eventTime,
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