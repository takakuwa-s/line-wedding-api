package driver

import (
	"github.com/gin-contrib/cors"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	wlc *controller.WeddingLineController
	alc *controller.AdminLineController
	uac *controller.UserApiController
	fac *controller.FileApiController
}

// Newコンストラクタ
func NewRouter(
	wlc *controller.WeddingLineController,
	alc *controller.AdminLineController,
	uac *controller.UserApiController,
	fac *controller.FileApiController) *Router {
	return &Router{wlc: wlc, alc:alc, uac:uac, fac:fac}
}

// Init ルーティング設定
func (r *Router) Init() {
	router := gin.Default()
	config := cors.DefaultConfig()
	frontDomain := os.Getenv("FRONT_DOMAIN")
	config.AllowOrigins = []string{frontDomain, "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(gin.Logger(), cors.New(config))
	router.POST("/line-messaging-api/wedding/webhook", r.wlc.Webhook)
	router.POST("/line-messaging-api/admin/webhook", r.alc.Webhook)
	router.PUT("/api/user", r.uac.UpdateUser)
	router.GET("/api/user/:id", r.uac.GetUser)
	router.GET("/api/file/list", r.fac.GetFileList)
	router.DELETE("/api/file/:id", r.fac.DeleteFile)
  port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}