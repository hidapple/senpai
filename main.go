package main

import (
	"fmt"
	"os"
)

var (
	token string
)

func main() {
	if token = os.Getenv("SLACK_API_TOKEN"); token == "" {
		fmt.Println("SLACK_API_TOKEN is empty.")
		os.Exit(1)
	}

	bot := NewBot(token)
	os.Exit(bot.Run())
}
