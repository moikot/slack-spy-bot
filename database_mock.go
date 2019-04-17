package main

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type UserCollectionMock struct {
	mock.Mock
}

func (m *UserCollectionMock) Set(ctx context.Context, user User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserCollectionMock) Get(ctx context.Context, userID string) (*User, error) {
	args := m.Called(ctx, userID)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).(*User), nil
	}
	return nil, err
}

func (m *UserCollectionMock) GetAllIDs(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	err := args.Error(1)
	if err == nil {
		return args.Get(0).([]string), nil
	}
	return nil, err
}

func (m *UserCollectionMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
