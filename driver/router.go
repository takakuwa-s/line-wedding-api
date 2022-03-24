package driver

import (
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"
  "github.com/gin-gonic/gin"
)

type Router struct {
	wlc *controller.WeddingLineController
	alc *controller.AdminLineController
}

// Newコンストラクタ
func NewRouter(wlc *controller.WeddingLineController, alc *controller.AdminLineController) *Router {
	return &Router{wlc: wlc, alc:alc}
}

// Init ルーティング設定
func (r *Router) Init() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.POST("/line-messaging-api/wedding/webhook", r.wlc.Webhook)
	router.POST("/line-messaging-api/admin/webhook", r.alc.Webhook)
  port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}