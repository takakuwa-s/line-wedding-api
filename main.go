package main

import (
	"log"
	"net/http"
	"os"
  "github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
  "github.com/gin-gonic/gin"
)

func webhook(c *gin.Context) {
	accessToken := os.Getenv("ACCESS_TOKEN")
	channelSecret := os.Getenv("CHANNEL_SECRET")

	client := &http.Client{}
	bot, err := linebot.New(channelSecret, accessToken, linebot.WithHTTPClient(client))
	if err != nil {
    c.String(http.StatusInternalServerError, "System error when creating client; err = %v", err)
		return
	}
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
    c.String(http.StatusInternalServerError, "System error when parsing the request; err = %v", err)
		return
	}
  log.Printf("Successfully parse the request; events = %v", events)
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			replyToken := event.ReplyToken
			leftBtn := linebot.NewMessageAction("left", "left clicked")
			rightBtn := linebot.NewMessageAction("right", "right clicked")
			template := linebot.NewConfirmTemplate("Hello World", leftBtn, rightBtn)
			message := linebot.NewTemplateMessage("Sorry :(, please update your app.", template)
			_, err := bot.ReplyMessage(replyToken, message).Do()
			if err != nil {
        c.String(http.StatusInternalServerError, "System error when replying the message; err = %v", err)
				return
			}
		}
	}
}

func main() {
  err := godotenv.Load("environments/dev.env")
	if err != nil {
		panic("Error loading .env file")
	}
  router := gin.Default()
  router.POST("/line-messaging-api/webhook", webhook)
  port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}
