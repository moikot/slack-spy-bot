package main

import (
	"testing"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Hello_SendsSubscribeUserPresence(t *testing.T) {
	users := &UserCollectionMock{}
	messenger := &MessengerMock{}

	bot := NewBot(users, messenger)

	ids := []string{"1", "2"}
	users.On("GetAllIDs", mock.Anything).Return(ids, nil)
	messenger.On("SubscribeUserPresence", ids).Return(nil)

	err := bot.Hello(&slack.HelloEvent{})

	assert.NoError(t, err)
	users.AssertExpectations(t)
	messenger.AssertExpectations(t)
}

func Test_PresenceChange_WhenPresenceChanges_SendsMessage(t *testing.T) {
	users := &UserCollectionMock{}
	messenger := &MessengerMock{}

	bot := NewBot(users, messenger)

	ev := &slack.PresenceChangeEvent{
		User:     "user",
		Presence: "presence",
	}

	user := &User{
		LastPresenceState:   "",
		NotificationChannel: "channel",
	}

	users.On("Get", mock.Anything, "user").Return(user, nil)
	users.On("Set", mock.Anything,
		mock.MatchedBy(func(user User) bool {
			return user.LastPresenceState == "presence"
		})).Return(nil)
	messenger.On("SendMessage", mock.Anything, "channel")

	err := bot.PresenceChange(ev)

	assert.NoError(t, err)
	users.AssertExpectations(t)
	messenger.AssertExpectations(t)
}

func Test_SpyOn_SetsTheUser_Notifies_Resubscribes(t *testing.T) {
	users := &UserCollectionMock{}
	messenger := &MessengerMock{}

	bot := NewBot(users, messenger)

	ev := &slack.MessageEvent{
		Msg: slack.Msg{
			Channel: "channel",
		},
	}

	users.On("Set", mock.Anything,
		mock.MatchedBy(func(user User) bool {
			return user.UserID == "user" &&
				user.NotificationChannel == "channel" &&
				user.LastPresenceState == "presence"
		})).Return(nil)

	ids := []string{"user"}
	users.On("GetAllIDs", mock.Anything).Return(ids, nil)

	presence := "presence"
	messenger.On("GetUserPresence", "user").Return(&presence, nil)
	messenger.On("SendMessage", mock.Anything, "channel")
	messenger.On("SubscribeUserPresence", ids).Return(nil)

	err := bot.SpyOn(ev, "user")

	assert.NoError(t, err)
	users.AssertExpectations(t)
	messenger.AssertExpectations(t)
}

func Test_SpyOn_DeletesTheUser_Notifies_Resubscribes(t *testing.T) {
	users := &UserCollectionMock{}
	messenger := &MessengerMock{}

	bot := NewBot(users, messenger)

	ev := &slack.MessageEvent{
		Msg: slack.Msg{
			Channel: "channel",
		},
	}

	user := &User{
		NotificationChannel: "channel",
	}
	users.On("Get", mock.Anything, "user").Return(user, nil)
	users.On("Delete", mock.Anything, "user").Return(nil)

	ids := []string{"user"}
	users.On("GetAllIDs", mock.Anything).Return(ids, nil)

	messenger.On("SendMessage", mock.Anything, "channel")
	messenger.On("SubscribeUserPresence", ids).Return(nil)

	err := bot.SpyOff(ev, "user")

	assert.NoError(t, err)
	users.AssertExpectations(t)
	messenger.AssertExpectations(t)
}
