package driver

import (
	"github.com/takakuwa-s/line-wedding-api/interface/controller"

	"os"
)

type WeddingRouter struct {
	cr  *CommonRouter
	lbc *controller.LineBotController
	iac *controller.InitApiController
	uac *controller.UserApiController
	fac *controller.FileApiController
	lac *controller.LineApiController
}

// Newコンストラクタ
func NewWeddingRouter(
	cr *CommonRouter,
	lbc *controller.LineBotController,
	iac *controller.InitApiController,
	uac *controller.UserApiController,
	fac *controller.FileApiController,
	lac *controller.LineApiController) *WeddingRouter {
	return &WeddingRouter{cr: cr, lbc: lbc, iac: iac, uac: uac, fac: fac, lac: lac}
}

// Init ルーティング設定
func (wr *WeddingRouter) Init() {
	router := wr.cr.GetDefaultRouter()

	router.POST("/line-messaging-api/wedding/webhook", wr.lbc.Webhook)
	router.GET("/health-check", wr.cr.HealthCheck)
	api := router.Group("/api")
	api.Use(wr.cr.ValidateTokenMiddleware)
	{
		user := api.Group("/user")
		{
			user.GET("/:id", wr.uac.GetUser)
			user.GET("/list", wr.uac.GetUserList)
			user.PUT("/:id", wr.uac.UpdateUser)
			user.PATCH("/:id", wr.uac.PatchUser)
		}
		file := api.Group("/file")
		{
			file.GET("/list", wr.fac.GetFileList)
			file.DELETE("/:id", wr.fac.DeleteFile)
			file.DELETE("/list", wr.fac.DeleteFileList)
		}
		api.GET("/init/:id", wr.iac.GetInitialData)
		api.POST("/line/message", wr.lac.SendMessageToLineBot)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
