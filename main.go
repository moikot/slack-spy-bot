package main

import (
	"context"
	"log"
	"os"
)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	ctx := context.Background()
	users, err := NewUserCollection(ctx)
	if err != nil {
		log.Fatal(err)
	}

	token := getenv("SPY_BOT_TOKEN")
	messenger := NewMessenger(token)
	bot := NewBot(users, messenger)

	messenger.AddHandler(bot.Hello)
	messenger.AddHandler(bot.PresenceChange)
	messenger.AddHandler(bot.Message)

	messenger.Listen()
}
