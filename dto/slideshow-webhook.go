package dto

type SlideshowWebhook struct {
	Type      string `json:"type"`
	Action    string `json:"action"`
	Id        string `json:"id"`
	Owner     string `json:"owner"`
	Status    string `json:"status"`
	Url       string `json:"url"`
	Error     string `json:"error"`
	Completed string `json:"completed"`
}
