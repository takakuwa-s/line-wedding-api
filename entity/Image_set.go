package entity

type ImageSet struct {
	Id      string   `json:"id"`
	Total   int      `json:"total"`
	FileIds []string `json:"fileIds"`
}

func NewImageSet(id string, total int) *ImageSet {
	return &ImageSet{
		Id:    id,
		Total: total,
	}
}
