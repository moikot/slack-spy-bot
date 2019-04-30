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

const users = "Users"

// A UserCollection provides access to users.
type UserCollection interface {
	Set(ctx context.Context, user User) error
	Get(ctx context.Context, userID string) (*User, error)
	GetAllIDs(ctx context.Context) ([]string, error)
	Delete(ctx context.Context, userID string) error
}

type fsDocument interface {
	Get(ctx context.Context) (*firestore.DocumentSnapshot, error)
	Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) (*firestore.WriteResult, error)
	Delete(ctx context.Context, preconds ...firestore.Precondition) (*firestore.WriteResult, error)
}

type fsDocumentImpl struct {
	*firestore.DocumentRef
}

type fsDocIterator interface {
	GetAll() ([]*firestore.DocumentSnapshot, error)
}

type fsDocIteratorImpl struct {
	*firestore.DocumentIterator
}

type fsCollection interface {
	Doc(id string) fsDocument
	Documents(ctx context.Context) fsDocIterator
}

type fsCollectionImpl struct {
	*firestore.CollectionRef
}

func (c *fsCollectionImpl) Doc(id string) fsDocument {
	return &fsDocumentImpl{DocumentRef: c.CollectionRef.Doc(id)}
}

func (c *fsCollectionImpl) Documents(ctx context.Context) fsDocIterator {
	return &fsDocIteratorImpl{DocumentIterator: c.CollectionRef.Documents(ctx)}
}

type fsClient interface {
	Collection(path string) fsCollection
}

type fsClientImpl struct {
	*firestore.Client
}

func (c *fsClientImpl) Collection(path string) fsCollection {
	return &fsCollectionImpl{CollectionRef: c.Client.Collection(path)}
}

type fsUserCollection struct {
	client fsClient
}

func NewUserCollection(ctx context.Context) (UserCollection, error) {
	client, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		return nil, err
	}

	c := &fsUserCollection{
		client: &fsClientImpl{Client: client},
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

func (fs *fsUserCollection) GetAllIDs(ctx context.Context) ([]string, error) {
	docs, err := fs.client.Collection(users).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var ids []string
	for _, doc := range docs {
		ids = append(ids, doc.Ref.ID)
	}

	return ids, nil
}

func (fs *fsUserCollection) Delete(ctx context.Context, userID string) error {
	_, err := fs.client.Collection(users).
		Doc(userID).
		Delete(ctx)

	return err
}
