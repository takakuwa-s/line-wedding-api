package driver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	frontDomain := os.Getenv("FRONT_DOMAIN")
	//TOTO localhost
	config.AllowOrigins = []string{frontDomain, "http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	router.Use(gin.Logger(), cors.New(config))
	return router
}

func (cr *CommonRouter) ValidateTokenMiddleware(c *gin.Context) {
	if err := cr.validateToken(c.GetHeader("Authorization")); err != nil {
		conf.Log.Error("Authorization failed", zap.String("error", err.Error()))
		//TOTO
		// c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		// c.Abort()
	}
}

func (cr *CommonRouter) validateToken(token string) error {
	client := &http.Client{}
	url := os.Getenv("LINE_API_BASE_URL")
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

func (cr *CommonRouter) HealthCheck(c *gin.Context) {
	c.Status(http.StatusOK)
}
