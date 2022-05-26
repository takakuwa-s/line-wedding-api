package driver

import (
	"github.com/gin-contrib/cors"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	lbc *controller.LineBotController
	iac *controller.InitApiController
	uac *controller.UserApiController
	fac *controller.FileApiController
	lac *controller.LineApiController
}

// Newコンストラクタ
func NewRouter(
	lbc *controller.LineBotController,
	iac *controller.InitApiController,
	uac *controller.UserApiController,
	fac *controller.FileApiController,
	lac *controller.LineApiController) *Router {
	return &Router{lbc: lbc, iac:iac, uac:uac, fac:fac, lac:lac}
}

// Init ルーティング設定
func (r *Router) Init() {
	router := gin.Default()
	config := cors.DefaultConfig()
	frontDomain := os.Getenv("FRONT_DOMAIN")
	config.AllowOrigins = []string{frontDomain, "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(gin.Logger(), cors.New(config))
	router.POST("/line-messaging-api/wedding/webhook", r.lbc.Webhook)
	router.PUT("/api/user", r.uac.UpdateUser)
	router.GET("/api/init/:id", r.iac.GetInitialData)
	router.GET("/api/file/list", r.fac.GetFileList)
	router.DELETE("/api/file/:id", r.fac.DeleteFile)
	router.POST("/api/line/message", r.lac.SendMessageToLineBot)
  port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}