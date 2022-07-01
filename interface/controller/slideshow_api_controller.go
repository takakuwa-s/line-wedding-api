package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type SlideShowApiController struct {
	su *usecase.SlideShowUsecase
}

// コンストラクタ
func NewSlideShowApiController(su *usecase.SlideShowUsecase) *SlideShowApiController {
	return &SlideShowApiController{su: su}
}

func (sac *SlideShowApiController) CreateSlideShow(c *gin.Context) {
	res, err := sac.su.CreateSlideShow()
	if err != nil {
		conf.Log.Error("[CreateSlideShow] Creating slide show failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"api": res})
}

func (suc *SlideShowApiController) UploadSlideShowWebhook(c *gin.Context) {
	var req dto.SlideshowWebhook
	if err := c.ShouldBindJSON(&req); err != nil {
		conf.Log.Error("[UploadSlideShow] json convestion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	conf.Log.Info("[UploadSlideShow] successfully received request", zap.Any("req", req))
	if err := suc.su.UploadSlideshow(req); err != nil {
		conf.Log.Error("[UploadSlideShow] Uploading slideshow failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (suc *SlideShowApiController) ListSlideshow(c *gin.Context) {
	list, err := suc.su.ListSlideshow()
	if err != nil {
		conf.Log.Error("[GetSlideShow] Getting slideshow failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"slideshow": list})
}

func (suc *SlideShowApiController) DeleteSlideshow(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := suc.su.DeleteSlideshow(id); err != nil {
		conf.Log.Error("[DeleteSlideshow] Deleting slideshow failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (suc *SlideShowApiController) PatchSlideshow(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	var req dto.PatchSlideShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		conf.Log.Error("[PatchSlideshow] json convestion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := suc.su.PatchSlideshow(id, req); err != nil {
		conf.Log.Error("[PatchSlideshow] Updating slideshow failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
