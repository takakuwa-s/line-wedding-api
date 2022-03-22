package dto

type File struct {
	MessageId string
}

func NewFile(messageId string) *File {
	return &File{
		MessageId: messageId,
	}
}