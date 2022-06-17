package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

type FileUploadController struct {
	fuu *usecase.FileUploadUsecase
}

// コンストラクタ
func NewFileUploadController(fuu *usecase.FileUploadUsecase) *FileUploadController {
	return &FileUploadController{fuu: fuu}
}

func (fuc *FileUploadController) UploadFile(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	c.Status(http.StatusAccepted)
	go func() {
		fuc.fuu.UploadFiles(ids)
	}()
}
