package dto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

type LineBot struct {
	client *linebot.Client
	token  string
	exp    time.Time
}

func NewLineBot() *LineBot {
	return &LineBot{}
}

func (lb *LineBot) GetClient() (*linebot.Client, error) {
	if lb.isExpired() {
		if err := lb.refreshToken(); err != nil {
			return nil, err
		}
		channelSecret := os.Getenv("LINE_BOT_CHANNEL_SECRET")
		bot, err := linebot.New(channelSecret, lb.token)
		if err != nil {
			return nil, fmt.Errorf("failed to create the wedding lineBot instance; err = %w", err)
		}
		lb.client = bot
	}
	return lb.client, nil
}

func (lb *LineBot) GetToken() (string, error) {
	if lb.isExpired() {
		if err := lb.refreshToken(); err != nil {
			return "", err
		}
	}
	return lb.token, nil
}

func (lb *LineBot) isExpired() bool {
	return time.Now().After(lb.exp)
}

func (lb *LineBot) refreshToken() error {
	res, err := lb.fetchToken()
	if err != nil {
		return err
	}
	lb.token = res["access_token"].(string)
	lb.exp = time.Now().Add(time.Second * time.Duration(res["expires_in"].(float64)-60*5))
	return nil
}

func (lb *LineBot) fetchToken() (map[string]interface{}, error) {
	jwt, err := lb.createJwt()
	if err != nil {
		return nil, err
	}
	form := url.Values{}
	form.Add("client_assertion", jwt)
	form.Add("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	client := &http.Client{}
	path := os.Getenv("LINE_API_BASE_URL") + "/oauth2/v2.1/token"
	resp, err := client.Post(path, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to call the line token generation api; %w", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the line token generation api response body; %w", err)
	}
	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("failed to json Unmarshal for the face line token generation api response; %w", err)
	}
	conf.Log.Info("Successfully get the access token", zap.Any("res", res))
	return res, nil
}

func (lb *LineBot) createJwt() (string, error) {
	channelId := os.Getenv("LINE_BOT_CHANNEL_ID")
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":       channelId,
		"sub":       channelId,
		"aud":       "https://api.line.me/",
		"exp":       time.Now().Add(time.Minute * 5).Unix(),
		"token_exp": 60 * 60 * 24,
	})
	token.Header["kid"] = os.Getenv("LINE_BOT_KID")

	path := os.Getenv("LINE_BOT_PRIVATE_KEY_PATH")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read the line bot private key; path = %s, err = %v", path, err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		return "", fmt.Errorf("failed to parses the line bot private key; path = %s, err = %v", path, err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt; err = %v", err)
	}
	return tokenString, nil
}
