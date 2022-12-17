package controller

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/dto"
)

type ConfigApiController struct {
}

// コンストラクタ
func NewConfigApiController() *ConfigApiController {
	return &ConfigApiController{}
}

func (cac *ConfigApiController) GetConfig(c *gin.Context) {
	fileFeatureAvailable, _ := strconv.ParseBool(os.Getenv("FILE_FEATURE_AVAILABLE"))
	attendanceFeatureAvailable, _ := strconv.ParseBool(os.Getenv("ATTENDANCE_FEATURE_AVAILABLE"))
	responce := dto.NewConfigResponce(fileFeatureAvailable, attendanceFeatureAvailable)
	c.JSON(http.StatusOK, gin.H{"config": responce})
}
