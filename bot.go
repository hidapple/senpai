package main

import (
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

const (
	ExitCodeOk    int = iota
	ExitCodeError int = iota

	IdeaFile string = "./idealist.txt"
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
			return ExitCodeError
		}
	}
	return ExitCodeOk
}

func add(bot *BOT, event *slack.MessageEvent) {
	file, err := os.OpenFile(IdeaFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Bye!")
	}
	defer file.Close()

	contents := strings.SplitAfter(event.Text, "<@"+botID+"> add")
	file.Write(([]byte)("-" + contents[1] + "\n"))
	bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage("added:+1:", event.Channel))
}

func show(bot *BOT, event *slack.MessageEvent) {
	if _, err := os.Stat(IdeaFile); err != nil {
		return
	}
	file, err := os.Open(IdeaFile)
	if err != nil {
		log.Fatal("Bye!")
	}
	defer file.Close()

	buf := make([]byte, 1024)
	message := ""
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			break
		}
		message += "\n" + string(buf[:n])
	}
	bot.rtm.SendMessage(bot.rtm.NewOutgoingMessage(message, event.Channel))
}
