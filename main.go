package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	token := getenv("SPY_BOT_TOKEN")
	userId := getenv("SPY_BOT_USER_ID")
	chanId := getenv("SPY_BOT_CHAN_ID")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Printf("Event received: %s\n", msg.Type)
			switch ev := msg.Data.(type) {

			case *slack.ConnectedEvent:
				ids := []string{userId};
				omsg := rtm.NewSubscribeUserPresence(ids);
				rtm.SendMessage(omsg)

			case *slack.PresenceChangeEvent:
				m := fmt.Sprintf("User: %v Presence: %v", ev.User, ev.Presence);
				omsg := rtm.NewOutgoingMessage(m, chanId);
				rtm.SendMessage(omsg)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				break Loop

			default:
			}
		}
	}
}
