package main

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/firestore"
)

// A User represents an observable user.
type User struct {
	UserID              string
	LastPresenceState   string
	NotificationChannel string
}

// A UserIterator allows to iterate a collection of users.
type UserIterator interface {
	Next() (*User, error)
	Stop()
}

type fsUserIterator struct {
	iter *firestore.DocumentIterator
}

func newUserIterator(iter *firestore.DocumentIterator) UserIterator {
	return &fsUserIterator{
		iter: iter,
	}
}

func (fs *fsUserIterator) Next() (*User, error) {
	doc, err := fs.iter.Next()
	if err != nil {
		return nil, err
	}

	user := &User{}

	err = doc.DataTo(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (fs *fsUserIterator) Stop() {
	fs.iter.Stop()
}

const users = "Users"

// A UserCollection provides access to users.
type UserCollection interface {
	Set(ctx context.Context, user User) error
	Get(ctx context.Context, userID string) (*User, error)
	GetAll(ctx context.Context) UserIterator
	Delete(ctx context.Context, userID string) error
}

type fsUserCollection struct {
	client *firestore.Client
}

func NewUserCollection(ctx context.Context) (UserCollection, error) {
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		return nil, err
	}

	c := &fsUserCollection{
		client: client,
	}
	return c, nil
}

func (fs *fsUserCollection) Set(ctx context.Context, user User) error {
	_, err := fs.client.Collection(users).
		Doc(user.UserID).
		Set(ctx, user)

	return err
}

func (fs *fsUserCollection) Get(ctx context.Context, userID string) (*User, error) {
	doc, err := fs.client.Collection(users).
		Doc(userID).
		Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	user := &User{}

	err = doc.DataTo(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (fs *fsUserCollection) GetAll(ctx context.Context) UserIterator {
	iter := fs.client.Collection(users).Documents(ctx)
	return newUserIterator(iter)
}

func (fs *fsUserCollection) Delete(ctx context.Context, userID string) error {
	_, err := fs.client.Collection(users).
		Doc(userID).
		Delete(ctx)

	return err
}
