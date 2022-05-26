package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/usecase/usecase"
	"go.uber.org/zap"
)

type FileApiController struct {
	au *usecase.ApiUsecase
}

// コンストラクタ
func NewFileApiController(au *usecase.ApiUsecase) *FileApiController {
	return &FileApiController{au: au}
}

func (fac *FileApiController) GetFileList(c *gin.Context) {
	err := fac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[GetFileList] Authorization failed", zap.String("error", err.Error()))
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "limit is required and must be number"})
		return
	}
	startId := c.Query("startId")
	userId := c.Query("userId")
	orderBy := c.Query("orderBy")
	files, err := fac.au.GetFileList(limit, startId, userId, orderBy)
	if err != nil {
		conf.Log.Error("[GetFileList] Getting file list failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (fac *FileApiController) DeleteFile(c *gin.Context) {
	err := fac.au.ValidateToken(c.GetHeader("Authorization"))
	if err != nil {
		conf.Log.Error("[DeleteFile] Authorization failed", zap.String("error", err.Error()))
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// return
	}
	id := c.Param("id")

	if err := fac.au.DeleteFile(id); err != nil {
		conf.Log.Error("[DeleteFile] Deleting file failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}