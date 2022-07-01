package ipresenter

type IPresenter interface {
	MulticastMessage(userIds []string, ms []map[string]interface{}) error
	ReplyMessage(token string, ms []map[string]interface{}) error
}
