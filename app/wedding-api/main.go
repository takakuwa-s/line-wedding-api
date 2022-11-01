package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/takakuwa-s/line-wedding-api/conf"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = os.Args[1]
	}
	err := godotenv.Load(fmt.Sprintf("../../env/%s/%s.env", env, env))
	if err != nil {
		conf.Log.Error("Error loading .env file", zap.Error(err))
		return
	}
	router := InitializeRouter()
	router.Init()
}
