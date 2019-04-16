package main

import (
	"regexp"
	"testing"

	"github.com/nlopes/slack"
)

func Test_RouteMessage_WhenRegExpMatches_CallsMessageHandler(t *testing.T) {
	router := NewRouter()
	var userId string

	regex := regexp.MustCompile(`foo\s*<@(\w+)>`)
	router.AddMessageHandler(regex, func(ev *slack.MessageEvent, id string) {
		userId = id
	})

	msg := &slack.MessageEvent{}
	msg.Text = "foo <@ID>"

	router.RouteMessage(msg)

	if userId != "ID" {
		t.Errorf("userId is '%v'; want 'ID'", userId)
	}
}

func Test_RouteEvent_WhenHandlerExists_CallsEventHandler(t *testing.T) {
	router := NewRouter()
	receivedEvent := false

	router.AddEventHandler(func(ev *slack.HelloEvent) {
		receivedEvent = true
	})

	router.RouteEvent(&slack.HelloEvent{})

	if !receivedEvent {
		t.Errorf("receivedEvent is 'false'; want 'true'")
	}
}

func Test_RouteEvent_WhenHandlerDoesNotExist_Succeeds(t *testing.T) {
	router := NewRouter()

	router.RouteEvent(&slack.HelloEvent{})
}
