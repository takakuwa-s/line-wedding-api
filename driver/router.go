package driver

import (
	"github.com/gin-contrib/cors"
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"

	"github.com/gin-gonic/gin"
)

type Router struct {
	lbc *controller.LineBotController
	cac *controller.CommonApiController
	iac *controller.InitApiController
	uac *controller.UserApiController
	fac *controller.FileApiController
	lac *controller.LineApiController
}

// Newコンストラクタ
func NewRouter(
	lbc *controller.LineBotController,
	cac *controller.CommonApiController,
	iac *controller.InitApiController,
	uac *controller.UserApiController,
	fac *controller.FileApiController,
	lac *controller.LineApiController) *Router {
	return &Router{lbc: lbc, cac: cac, iac: iac, uac: uac, fac: fac, lac: lac}
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
	api := router.Group("/api")
	api.Use(r.cac.ValidateTokenMiddleware)
	{
		user := api.Group("/user")
		{
			user.GET("/:id", r.uac.GetUser)
			user.GET("/list", r.uac.GetUserList)
			user.PUT("/:id", r.uac.UpdateUser)
			user.PATCH("/:id", r.uac.PatchUser)
		}
		file := api.Group("/file")
		{
			file.GET("/list", r.fac.GetFileList)
			file.DELETE("/:id", r.fac.DeleteFile)
			file.DELETE("/list", r.fac.DeleteFileList)
		}
		api.GET("/init/:id", r.iac.GetInitialData)
		api.POST("/line/message", r.lac.SendMessageToLineBot)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
