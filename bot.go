package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/nlopes/slack"
	"google.golang.org/api/iterator"
)

type Bot struct {
	users     UserCollection
	messenger Messenger
}

var (
	SpyOnRegEx  = regexp.MustCompile(`(?i)spy-on\s+<@(\w+)>`)
	SpyOffRegEx = regexp.MustCompile(`(?i)spy-off\s+<@(\w+)>`)
)

func NewBot(users UserCollection, messenger Messenger) *Bot {
	return &Bot{
		users:     users,
		messenger: messenger,
	}
}

func (b *Bot) Hello(ev *slack.HelloEvent) error {
	return b.resubscribe()
}

func (b *Bot) resubscribe() error {
	ctx := context.Background()

	var ids []string
	iter := b.users.GetAll(ctx)
	defer iter.Stop()
	for {
		user, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		ids = append(ids, user.UserID)
	}

	if len(ids) > 0 {
		b.messenger.SubscribeUserPresence(ids)
	}

	return nil
}

func (b *Bot) PresenceChange(ev *slack.PresenceChangeEvent) error {
	ctx := context.Background()
	user, err := b.users.Get(ctx, ev.User)
	if err != nil || user == nil {
		return err
	}

	if user.LastPresenceState != ev.Presence {
		b.notifyPresenceChanged(ev.User, ev.Presence, user.NotificationChannel)

		// Save new presence
		user.LastPresenceState = ev.Presence
		return b.users.Set(ctx, *user)
	}

	return nil
}

func (b *Bot) SpyOn(ev *slack.MessageEvent, userID string) error {
	presence, err := b.messenger.GetUserPresence(userID)
	if err != nil {
		return err
	}

	user := User{
		UserID:              userID,
		NotificationChannel: ev.Channel,
		LastPresenceState:   *presence,
	}

	ctx := context.Background()
	err = b.users.Set(ctx, user)
	if err != nil {
		return err
	}

	b.notifyPresenceChanged(userID, *presence, ev.Channel)
	return b.resubscribe()
}

func (b *Bot) SpyOff(ev *slack.MessageEvent, userID string) error {
	ctx := context.Background()
	user, err := b.users.Get(ctx, userID)
	if err != nil {
		return err
	}

	err = b.users.Delete(ctx, userID)
	if err != nil {
		return err
	}

	b.notifyStoppedSpying(userID, user.NotificationChannel)
	return b.resubscribe()
}

func (b *Bot) notifyPresenceChanged(userID string, presence string, channelID string) {
	msg := fmt.Sprintf("User: <@%v> Presence: %v", userID, presence)
	b.messenger.SendMessage(msg, channelID)
}

func (b *Bot) notifyStoppedSpying(userID string, channelID string) {
	msg := fmt.Sprintf("Stopped spying on: <@%v>", userID)
	b.messenger.SendMessage(msg, channelID)
}
