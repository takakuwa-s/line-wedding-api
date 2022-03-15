package driver

import (
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"
  "github.com/gin-gonic/gin"
)

type Router struct {
	lc *controller.LineController
}

// Newコンストラクタ
func NewRouter(lc *controller.LineController) *Router {
	return &Router{lc: lc}
}

// Init ルーティング設定
func (r *Router) Init() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.POST("/line-messaging-api/webhook", r.lc.Webhook)
  port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}