package main

import (
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("environments/dev.env")
	if err != nil {
		panic("Error loading .env file")
	}
	router := InitializeRouter()
	router.Init()
}
