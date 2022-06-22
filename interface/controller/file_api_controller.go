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
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit is required and must be number"})
		return
	}
	needCreaterName, err := strconv.ParseBool(c.Query("needCreaterName"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "needCreaterName must be boolean"})
		return
	}
	uploadedStr := c.Query("uploaded")
	var uploaded *bool
	if uploadedStr != "" {
		u, err := strconv.ParseBool(uploadedStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "uploaded must be boolean"})
			return
		}
		uploaded = &u
	}
	startId := c.Query("startId")
	userId := c.Query("userId")
	orderBy := c.Query("orderBy")
	files, err := fac.au.GetFileList(limit, startId, userId, orderBy, uploaded, needCreaterName)
	if err != nil {
		conf.Log.Error("[GetFileList] Getting file list failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (fac *FileApiController) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	if err := fac.au.DeleteFile(id); err != nil {
		conf.Log.Error("[DeleteFile] Deleting file failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (fac *FileApiController) DeleteFileList(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if err := fac.au.DeleteFileList(ids); err != nil {
		conf.Log.Error("[DeleteFileList] Deleting multiple files failed", zap.String("error", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
