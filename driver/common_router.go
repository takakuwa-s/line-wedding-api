package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type CommonRouter struct {
}

// コンストラクタ
func NewCommonRouter() *CommonRouter {
	return &CommonRouter{}
}

func (cr *CommonRouter) GetDefaultRouter() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	frontUrl := os.Getenv("FRONT_URL")
	config.AllowOrigins = []string{frontUrl, "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(gin.Logger(), cors.New(config))
	return router
}

func (cr *CommonRouter) ValidateTokenMiddleware(c *gin.Context, channelId string) {
	if os.Getenv("ENV") == "local" {
		return
	}

	auth := c.GetHeader("Authorization")
	idx := strings.Index(auth, "Bearer ")
	if idx == -1 || len(auth) <= 7 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer {access token} is required"})
		c.Abort()
		return
	}
	token := auth[idx+7:]
	if err := cr.validateToken(token, channelId); err != nil {
		conf.Log.Error("Authorization failed", zap.String("error", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
	}
}

func (cr *CommonRouter) validateToken(token, channelId string) error {
	client := &http.Client{}
	url := os.Getenv("LINE_API_BASE_URL")
	resp, err := client.Get(url + "/oauth2/v2.1/verify?access_token=" + token)
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
	} else if obj["client_id"].(string) != channelId {
		return fmt.Errorf("invalid access token audience; client_id = %s", obj["client_id"].(string))
	} else if obj["expires_in"].(float64) <= 0 {
		return fmt.Errorf("access token expired; expires_in=%f", obj["expires_in"].(float64))
	}
	conf.Log.Info("Successfully validate the token")
	return nil
}

func (cr *CommonRouter) HealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
