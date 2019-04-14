package main

import (
	"context"
	"encoding/json"
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
	spyOnRegEx  = regexp.MustCompile(`spy-on\s+<@(\w+)>`)
	spyOffRegEx = regexp.MustCompile(`spy-off\s+<@(\w+)>`)
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

		// Save new state
		user.LastPresenceState = ev.Presence
		return b.users.Set(ctx, *user)
	}

	return nil
}

func (b *Bot) Message(ev *slack.MessageEvent) error {
	s, _ := json.Marshal(ev)
	fmt.Printf("Message: %v\n", string(s))

	match := spyOnRegEx.FindStringSubmatch(ev.Text)
	if match != nil {
		return b.spyOn(ev.Channel, match[1])
	}

	match = spyOffRegEx.FindStringSubmatch(ev.Text)
	if match != nil {
		return b.spyOff(match[1])
	}

	return nil
}

func (b *Bot) spyOn(channelID string, userID string) error {
	presence, err := b.messenger.GetUserPresence(userID)
	if err != nil {
		return err
	}

	user := User{
		UserID:              userID,
		NotificationChannel: channelID,
		LastPresenceState:   *presence,
	}

	ctx := context.Background()
	err = b.users.Set(ctx, user)
	if err != nil {
		return err
	}

	b.notifyPresenceChanged(userID, *presence, channelID)
	return b.resubscribe()
}

func (b *Bot) spyOff(userID string) error {
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
