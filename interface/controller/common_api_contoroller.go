package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type CommonApiController struct {
}

// コンストラクタ
func NewCommonApiController() *CommonApiController {
	return &CommonApiController{}
}

func (cac *CommonApiController) ValidateTokenMiddleware(c *gin.Context) {
	if err := cac.validateToken(c.GetHeader("Authorization")); err != nil {
		conf.Log.Error("Authorization failed", zap.String("error", err.Error()))
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// c.Abort()
	}
}

func (cac *CommonApiController) validateToken(token string) error {
	client := &http.Client{}
	url := os.Getenv("LIFF_API_BASE_URL")
	resp, err := client.Get(url + "?access_token=" + token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("token validation failed; %s", obj["error_description"])
	} else if obj["client_id"].(string) != os.Getenv("LIFF_CHANNEL_ID") {
		return fmt.Errorf("invalid access token audience")
	} else if obj["expires_in"].(float64) <= 0 {
		return fmt.Errorf("access token expired")
	}
	return nil
}
