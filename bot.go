package main

import (
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	botID   string
	botName string

	commands = map[string]string{
		"ADD":  "add",
		"SHOW": "show",
	}
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
			if strings.HasPrefix(ev.Text, "<@"+botID+"> "+commands["ADD"]) {
				add(bot, ev)
			}

			if strings.HasPrefix(ev.Text, "<@"+botID+"> "+commands["SHOW"]) {
				show(bot, ev)
			}

		case *slack.InvalidAuthEvent:
			return 1
		}
	}
	return 0
}

func add(bot *BOT, event *slack.MessageEvent) {
	file, _ := os.OpenFile("./idealist.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer file.Close()

	contents := strings.SplitAfter(event.Text, "<@"+botID+"> add")
	file.Write(([]byte)("-" + contents[1] + "\n"))
	bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("added :+1:", event.Channel))
}

func show(bot *BOT, event *slack.MessageEvent) {
	file, _ := os.Open("./idealist.txt")
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
	bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(message, event.Channel))
}
