package dto

import (
	"encoding/json"
	"fmt"

	"github.com/takakuwa-s/line-wedding-api/entity"
)

type TextMessage struct {
	ReplyToken   string
	Text         string
	SenderUserId string
}

func NewTextMessage(replyToken, text, senderUserId string) *TextMessage {
	return &TextMessage{
		ReplyToken:   replyToken,
		Text:         text,
		SenderUserId: senderUserId,
	}
}

type ImageMessage struct {
	ReplyToken string
	File       *entity.File
	ImageSet   *entity.ImageSet
}

func NewImageMessage(replyToken string, file *entity.File, imageSet *entity.ImageSet) *ImageMessage {
	return &ImageMessage{
		ReplyToken: replyToken,
		File:       file,
		ImageSet:   imageSet,
	}
}

type VideoMessage struct {
	ReplyToken string
	File       *entity.File
	Duration   int
}

func NewVideoMessage(replyToken string, file *entity.File, duration int) *VideoMessage {
	return &VideoMessage{
		ReplyToken: replyToken,
		File:       file,
		Duration:   duration,
	}
}

type FollowMessage struct {
	ReplyToken   string
	SenderUserId string
}

func NewFollowMessage(replyToken, senderUserId string) *FollowMessage {
	return &FollowMessage{
		ReplyToken:   replyToken,
		SenderUserId: senderUserId,
	}
}

type PostbackMessage struct {
	ReplyToken   string
	Data         map[string]interface{}
	Params       interface{}
	SenderUserId string
}

func NewPostbackMessage(replyToken, dataStr, senderUserId string, params interface{}) *PostbackMessage {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		panic(fmt.Sprintf("Failed to convert postback data to map object; dataStr = %s, err = %v", dataStr, err))
	}
	return &PostbackMessage{
		ReplyToken:   replyToken,
		Data:         data,
		Params:       params,
		SenderUserId: senderUserId,
	}
}

type GroupMessage struct {
	ReplyToken string
}

func NewGroupMessage(replyToken string) *GroupMessage {
	return &GroupMessage{
		ReplyToken: replyToken,
	}
}
