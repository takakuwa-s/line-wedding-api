package dto

type MulticastMessage struct {
	UserIds []string
	Messages []map[string]interface{}
}

func NewMulticastMessage(userIds []string, messages []map[string]interface{}) *MulticastMessage {
	return &MulticastMessage{
		UserIds : userIds,
		Messages : messages,
	}
} 