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

	token := getenv("BOT_TOKEN")
	router := NewRouter()

	messenger := NewMessenger(token, router)
	bot := NewBot(users, messenger)

	router.AddEventHandler(bot.Hello)
	router.AddEventHandler(bot.PresenceChange)
	router.AddMessageHandler(SpyOnRegEx, bot.SpyOn)
	router.AddMessageHandler(SpyOffRegEx, bot.SpyOff)

	messenger.Listen()
}
