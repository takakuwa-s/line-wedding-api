package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type UserApiController struct {
	au *usecase.ApiUsecase
}

// コンストラクタ
func NewUserApiController(au *usecase.ApiUsecase) *UserApiController {
	return &UserApiController{au: au}
}

func (uac *UserApiController) UpdateUser(c *gin.Context) {
	err := uac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[UpdateUser] Authorization failed", zap.String("error", err.Error()))
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// return
	}
	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		conf.Log.Error("[UpdateUser] json convestion failed", zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := uac.au.UpdateUser(&request)
	if err != nil {
		conf.Log.Error("[UpdateUser] Updating user failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
