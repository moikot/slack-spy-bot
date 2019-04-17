package main

import (
	"regexp"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

type RouterMock struct {
	mock.Mock
}

func (m *RouterMock) RouteEvent(event interface{}) {
	m.Called(event)
}

func (m *RouterMock) RouteMessage(event *slack.MessageEvent) {
	m.Called(event)
}

func (m *RouterMock) AddEventHandler(handler HandlerFunc) {
	m.Called(handler)
}

func (m *RouterMock) AddMessageHandler(regexp *regexp.Regexp, handler HandlerFunc) {
	m.Called(regexp, handler)
}
