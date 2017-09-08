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

func main() {
	api := slack.New("")
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
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
					return
				}
				defer file.Close()

				contents := strings.SplitAfter(ev.Text, "<@"+botID+"> add")
				file.Write(([]byte)("-" + contents[1] + "\n"))
				rtm.SendMessage(rtm.NewOutgoingMessage("added :+1:", channel))
			}

			if strings.HasPrefix(text, "<@"+botID+"> show") {
				file, err := os.Open("./idealist.txt")
				if err != nil {
					return
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
				rtm.SendMessage(rtm.NewOutgoingMessage(message, channel))
			}

		case *slack.InvalidAuthEvent:
			return
		}
	}
}
