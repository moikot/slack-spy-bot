package main

import (
	"reflect"
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/mock"
)

type RTMMock struct {
	mock.Mock
}

func (m *RTMMock) GetUserPresence(user string) (*slack.UserPresence, error) {
	args := m.Called(user)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*slack.UserPresence), nil
	}
	return nil, err
}

func (m *RTMMock) SendMessage(msg *slack.OutgoingMessage) {
	m.Called(msg)
}

func (m *RTMMock) ManageConnection() {
	m.Called()
}

func Test_SubscribeUserPresence_SendsSubscribeUserPresence(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	ids := []string{"id1", "id2"}
	rtm.On("SendMessage",
		mock.MatchedBy(
			func(req *slack.OutgoingMessage) bool {
				return req.Type == "presence_sub" &&
					reflect.DeepEqual(req.IDs, ids)
			})).Once()

	m.SubscribeUserPresence(ids)

	rtm.AssertExpectations(t)
}

func Test_GetUserPresence_WhenSucceeds_ReturnsUserPresence(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	id := "id1"
	up := &slack.UserPresence{Presence: "presence"}

	rtm.On("GetUserPresence", id).Once().Return(up, nil)

	presence, _ := m.GetUserPresence(id)

	if *presence != "presence" {
		t.Errorf("user presence is '%v'; want 'presence'", *presence)
	}
}

func Test_SendMessage_SendsOutgoingMessage(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	text := "text"
	channel := "channel"

	rtm.On("SendMessage",
		mock.MatchedBy(
			func(req *slack.OutgoingMessage) bool {
				return req.Type == "message" &&
					req.Text == text &&
					req.Channel == channel
			})).Once()

	m.SendMessage(text, channel)

	rtm.AssertExpectations(t)
}

func Test_Listen_WhenMessageEventIsReceived_CallsRouteMessage(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	ev := &slack.MessageEvent{}
	rtm.On("ManageConnection").Run(func(args mock.Arguments) {
		m.rtm.IncomingEvents <- slack.RTMEvent{
			Data: ev,
		}
		// Exit from the loop
		m.rtm.IncomingEvents <- slack.RTMEvent{
			Data: &slack.DisconnectedEvent{Intentional: true},
		}
	})

	router := &RouterMock{}
	router.On("RouteMessage", ev)

	m.router = router
	m.Listen()

	router.AssertExpectations(t)
}

func Test_Listen_WhenHelloEventIsReceived_CallsRouteEvent(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	ev := &slack.HelloEvent{}
	rtm.On("ManageConnection").Run(func(args mock.Arguments) {
		m.rtm.IncomingEvents <- slack.RTMEvent{
			Data: ev,
		}
		// Exit from the loop
		m.rtm.IncomingEvents <- slack.RTMEvent{
			Data: &slack.DisconnectedEvent{Intentional: true},
		}
	})

	router := &RouterMock{}
	router.On("RouteEvent", ev)

	m.router = router
	m.Listen()

	router.AssertExpectations(t)
}

func Test_Listen_WhenInvalidAuthEventIsReceived_Exits(t *testing.T) {
	rtm := &RTMMock{}
	m := newStubbedMessenger(rtm)

	ev := &slack.InvalidAuthEvent{}
	rtm.On("ManageConnection").Run(func(args mock.Arguments) {
		m.rtm.IncomingEvents <- slack.RTMEvent{
			Data: ev,
		}
	})

	router := &RouterMock{}

	m.router = router
	m.Listen()

	router.AssertExpectations(t)
}

func newStubbedMessenger(rtm rtm_) *slackMessenger {
	api := slack.New("")

	return &slackMessenger{
		rtm:  api.NewRTM(),
		rtm_: rtm,
	}
}
