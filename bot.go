package main

import (
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	botID   string
	botName string
)

type BOT struct {
	api *slack.Client
	rtm *slack.RTM
}

func NewBot(token string) *BOT {
	bot := new(BOT)
	bot.api = slack.New(token)
	bot.rtm = bot.api.NewRTM()
	return bot
}

func (bot *BOT) Run() int {
	go bot.rtm.ManageConnection()

	for msg := range bot.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			botID = ev.Info.User.ID
			botName = ev.Info.User.Name

		case *slack.MessageEvent:
			text := ev.Text
			channel := ev.Channel

			if strings.HasPrefix(text, "<@"+botID+"> add") {
				file, err := os.OpenFile("./idealist.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					return 1
				}
				defer file.Close()

				contents := strings.SplitAfter(ev.Text, "<@"+botID+"> add")
				file.Write(([]byte)("-" + contents[1] + "\n"))
				bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("added :+1:", channel))
			}

			if strings.HasPrefix(text, "<@"+botID+"> show") {
				file, err := os.Open("./idealist.txt")
				if err != nil {
					return 1
				}
				defer file.Close()

				buf := make([]byte, 1024)
				message := ""
				for {
					n, err := file.Read(buf)
					if err != nil {
						break
					}
					if n == 0 {
						break
					}
					message += "\n" + string(buf[:n])
				}
				bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(message, channel))
			}

		case *slack.InvalidAuthEvent:
			return 1
		}
	}
	return 0
}
