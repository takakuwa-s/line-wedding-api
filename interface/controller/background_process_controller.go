package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
)

type BackgroundProcessController struct {
	bpu *usecase.BackgroundProcessUsecase
}

// コンストラクタ
func NewBackgroundProcessController(fuu *usecase.BackgroundProcessUsecase) *BackgroundProcessController {
	return &BackgroundProcessController{bpu: fuu}
}

func (bpc *BackgroundProcessController) UploadFile(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	c.Status(http.StatusAccepted)
	go func() {
		bpc.bpu.UploadFilesByIds(ids)
	}()
}

func (bpc *BackgroundProcessController) DeleteFileList(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	c.Status(http.StatusAccepted)
	go func() {
		bpc.bpu.DeleteFilesByIds(ids)
	}()
}
