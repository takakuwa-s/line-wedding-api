package driver

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"
)

type FileUploadRouter struct {
	cr  *CommonRouter
	fuc *controller.FileUploadController
}

// Newコンストラクタ
func NewFileUploadRouter(
	cr *CommonRouter,
	fuc *controller.FileUploadController) *FileUploadRouter {
	return &FileUploadRouter{cr: cr, fuc: fuc}
}

// Init ルーティング設定
func (fur *FileUploadRouter) Init() {
	router := fur.cr.GetDefaultRouter()
	router.Use(fur.validateTokenMiddleware)
	router.POST("/api/file/list", fur.fuc.UploadFile)
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
