package dto

import (
	"github.com/takakuwa-s/line-wedding-api/entity"
)

type FileResponce struct {
	entity.File
	CreaterName string `json:"createrName"`
}

func NewFileResponce(file entity.File, createrName string) FileResponce {
	return FileResponce{
		file,
		createrName,
	}
}

func NewFileResponceList(files []entity.File, uMap map[string]string) []FileResponce {
	var res []FileResponce
	for _, f := range files {
		res = append(res, NewFileResponce(f, uMap[f.Creater]))
	}
	return res
}
