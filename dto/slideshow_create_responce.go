package dto

type QueuedResponseData struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type SlideShowCreateResponce struct {
	Success  bool               `json:"success"`
	Message  string             `json:"message"`
	Response QueuedResponseData `json:"response"`
}
