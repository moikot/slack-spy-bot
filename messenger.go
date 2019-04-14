package main

import (
	"fmt"
	"reflect"

	"github.com/nlopes/slack"
)

type HandlerFunc interface{}

type Messenger interface {
	SubscribeUserPresence(ids []string)
	GetUserPresence(userID string) (*string, error)
	SendMessage(text string, channelID string)

	AddHandler(handler HandlerFunc)
	Listen()
}

func NewMessenger(token string) Messenger {
	api := slack.New(token)
	rtm := api.NewRTM()

	return &slackMessenger{
		handlers: map[string]HandlerFunc{},
		rtm:      rtm,
	}
}

type slackMessenger struct {
	handlers map[string]HandlerFunc
	rtm      *slack.RTM
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

func (m *slackMessenger) AddHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventType := handlerType.In(0).Elem().Name()
	m.handlers[eventType] = handler
}

func (m *slackMessenger) Listen() {
	go m.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-m.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.RTMError:
				fmt.Printf("error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("error: invalid authentication\n")
				break Loop

			default:
				m.routeEvent(ev)
			}
		}
	}
}

func (m *slackMessenger) routeEvent(event interface{}) {
	var eventType = reflect.TypeOf(event).Elem().Name()

	var handler = m.handlers[eventType]
	if handler == nil {
		return
	}

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(event)

	reflect.ValueOf(handler).Call(params)
}
