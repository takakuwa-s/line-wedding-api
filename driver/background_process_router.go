package driver

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
)

type BackgroundProcessRouter struct {
	cr  *CommonRouter
	bpc *controller.BackgroundProcessController
	sac *controller.SlideShowApiController
}

// Newコンストラクタ
func NewBackgroundProcessRouter(
	cr *CommonRouter,
	bpc *controller.BackgroundProcessController,
	sac *controller.SlideShowApiController) *BackgroundProcessRouter {
	return &BackgroundProcessRouter{cr: cr, bpc: bpc, sac: sac}
}

// Init ルーティング設定
func (bpr *BackgroundProcessRouter) Init() {
	router := bpr.cr.GetDefaultRouter()
	api := router.Group("/api")
	api.Use(bpr.validateTokenMiddleware)
	router.GET("/health-check", bpr.cr.HealthCheck)
	{
		api.POST("/file/list", bpr.bpc.UploadFile)
		api.DELETE("/file/list", bpr.bpc.DeleteFileList)
	}
	router.POST("/webhook/slideshow", bpr.sac.UploadSlideShowWebhook)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	router.Run(":" + port)
}

func (bpr *BackgroundProcessRouter) validateTokenMiddleware(c *gin.Context) {
	bpr.cr.ValidateTokenMiddleware(c, os.Getenv("LINE_BOT_CHANNEL_ID"))
}
