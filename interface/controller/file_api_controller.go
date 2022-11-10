package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"github.com/takakuwa-s/line-wedding-api/dto"
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
	forBrideAndGroomStr := c.Query("forBrideAndGroom")
	var forBrideAndGroom *bool
	if forBrideAndGroomStr != "" {
		f, err := strconv.ParseBool(forBrideAndGroomStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "forBrideAndGroom must be boolean"})
			return
		}
		forBrideAndGroom = &f
	}
	startId := c.Query("startId")
	userId := c.Query("userId")
	fileType := c.Query("fileType")
	orderBy := c.Query("orderBy")
	fileStatuses := c.QueryArray("fileStatus")
	files, err := fac.au.GetFileList(limit, startId, userId, orderBy, fileType, fileStatuses, forBrideAndGroom, needCreaterName)
	if err != nil {
		conf.Log.Error("[GetFileList] Getting file list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files})
}

func (fac *FileApiController) DeleteFileList(c *gin.Context) {
	ids, exists := c.GetQueryArray("id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if err := fac.au.DeleteFileList(ids); err != nil {
		conf.Log.Error("[DeleteFileList] Deleting multiple files failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (fac *FileApiController) PatchFile(c *gin.Context) {
	var request dto.PatchFileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		conf.Log.Error("[PatchFile] json convestion failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	if err := fac.au.PatchFile(id, request.ForBrideAndGroom); err != nil {
		conf.Log.Error("[PatchFile] Updating file field failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
