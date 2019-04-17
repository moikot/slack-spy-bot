package main

import (
	"github.com/stretchr/testify/mock"
)

type MessengerMock struct {
	mock.Mock
}

func (m *MessengerMock) SubscribeUserPresence(ids []string) {
	m.Called(ids)
}

func (m *MessengerMock) GetUserPresence(userID string) (*string, error) {
	args := m.Called(userID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*string), nil
	}
	return nil, err
}

func (m *MessengerMock) SendMessage(text string, channelID string) {
	m.Called(text, channelID)
}

func (m *MessengerMock) Listen() {
	m.Called()
}
