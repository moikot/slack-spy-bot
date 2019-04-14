package main

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/nlopes/slack"
)

type HandlerFunc interface{}

type Messenger interface {
	SubscribeUserPresence(ids []string)
	GetUserPresence(userID string) (*string, error)
	SendMessage(text string, channelID string)

	AddEventHandler(handler HandlerFunc)
	AddMessageHandler(regex *regexp.Regexp, handler HandlerFunc)

	Listen()
}

type slackMessenger struct {
	eventHandlers   map[string]HandlerFunc
	messageHandlers map[*regexp.Regexp]HandlerFunc
	rtm             *slack.RTM
}

func NewMessenger(token string) Messenger {
	api := slack.New(token)
	rtm := api.NewRTM()

	return &slackMessenger{
		eventHandlers:   map[string]HandlerFunc{},
		messageHandlers: map[*regexp.Regexp]HandlerFunc{},
		rtm:             rtm,
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

func (m *slackMessenger) AddEventHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventType := handlerType.In(0).Elem().Name()
	m.eventHandlers[eventType] = handler
}

func (m *slackMessenger) AddMessageHandler(regexp *regexp.Regexp, handler HandlerFunc) {
	m.messageHandlers[regexp] = handler
}

func (m *slackMessenger) Listen() {
	go m.rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-m.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				m.routeMessage(ev)

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

func (m *slackMessenger) routeMessage(event *slack.MessageEvent) {
	// Try to find matching pattern
	for pattern, handler := range m.messageHandlers {
		match := pattern.FindAllStringSubmatch(event.Text, -1)
		pattern.SubexpNames()
		if match != nil {
			for _, captures := range match {
				var params = make([]reflect.Value, len(captures))
				params[0] = reflect.ValueOf(event)

				for index, capture := range captures {
					if index == 0 {
						continue
					}
					params[index] = reflect.ValueOf(capture)
				}

				reflect.ValueOf(handler).Call(params)
			}
		}
	}
}

func (m *slackMessenger) routeEvent(event interface{}) {
	var eventType = reflect.TypeOf(event).Elem().Name()

	var handler = m.eventHandlers[eventType]
	if handler == nil {
		return
	}

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(event)

	reflect.ValueOf(handler).Call(params)
}
