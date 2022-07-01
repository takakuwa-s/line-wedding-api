package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type InitApiController struct {
	au *usecase.ApiUsecase
}

// コンストラクタ
func NewInitApiController(au *usecase.ApiUsecase) *InitApiController {
	return &InitApiController{au: au}
}

func (iac *InitApiController) GetInitialData(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	data, err := iac.au.GetInitialData(id)
	if err != nil {
		conf.Log.Error("[GetInitialData] Getting user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if data.User == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
