package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

type HandlerFunc interface{}

type Messenger interface {
	SubscribeUserPresence(ids []string)
	GetUserPresence(userID string) (*string, error)
	SendMessage(text string, channelID string)

	Listen()
}

type slackMessenger struct {
	rtm    *slack.RTM
	router Router
}

func NewMessenger(token string, router Router) Messenger {
	api := slack.New(token)
	rtm := api.NewRTM()

	return &slackMessenger{
		router: router,
		rtm:    rtm,
	}
}

func (m *slackMessenger) SubscribeUserPresence(ids []string) {
	msg := m.rtm.NewSubscribeUserPresence(ids)
	m.rtm.SendMessage(msg)
}

func (m *slackMessenger) GetUserPresence(userID string) (*string, error) {
	presence, err := m.rtm.GetUserPresence(userID)
	if err != nil {
		return nil, err
	}
	return &presence.Presence, nil
}

func (m *slackMessenger) SendMessage(text string, channelID string) {
	msg := m.rtm.NewOutgoingMessage(text, channelID)
	m.rtm.SendMessage(msg)
}

func (m *slackMessenger) Listen() {
	go m.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-m.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				m.router.RouteMessage(ev)

			case *slack.RTMError:
				fmt.Printf("error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("error: invalid authentication\n")
				break Loop

			default:
				m.router.RouteEvent(ev)
			}
		}
	}
}
