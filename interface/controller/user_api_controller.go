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

func (uac *UserApiController) GetUser(c *gin.Context) {
	err := uac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[GetUser] Authorization failed", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	user, err := uac.au.GetUser(id)
	if err != nil {
		conf.Log.Error("[GetUser] Getting user failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (uac *UserApiController) UpdateUser(c *gin.Context) {
	err := uac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[UpdateUser] Authorization failed", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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
