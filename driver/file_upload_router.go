package driver

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
)

type FileUploadRouter struct {
	cr  *CommonRouter
	fuc *controller.FileUploadController
	sac *controller.SlideShowApiController
}

// Newコンストラクタ
func NewFileUploadRouter(
	cr *CommonRouter,
	fuc *controller.FileUploadController,
	sac *controller.SlideShowApiController) *FileUploadRouter {
	return &FileUploadRouter{cr: cr, fuc: fuc, sac: sac}
}

// Init ルーティング設定
func (fur *FileUploadRouter) Init() {
	router := fur.cr.GetDefaultRouter()
	api := router.Group("/api")
	api.Use(fur.validateTokenMiddleware)
	{
		api.POST("/file/list", fur.fuc.UploadFile)
	}
	router.POST("/webhook/slideshow", fur.sac.UploadSlideShowWebhook)
	router.GET("/health-check", fur.cr.HealthCheck)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}

func (fur *FileUploadRouter) validateTokenMiddleware(c *gin.Context) {
	fur.cr.ValidateTokenMiddleware(c, os.Getenv("LINE_BOT_CHANNEL_ID"))
}
