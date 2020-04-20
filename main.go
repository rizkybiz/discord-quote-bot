package main

import (
	"flag"
	"log"

	"github.com/rizkybiz/discord-quote-bot/bot"
)

func main() {

	//Parse run flags
	var token string
	flag.StringVar(&token, "token", "", "token of the discord bot")
	flag.Parse()

	// Check if the token has been provided
	if token == "" {
		log.Fatalf("Discord bot token not set, exiting")
	}

	//Initialize the bot and connect
	bot, err := bot.New(token)
	if err != nil {
		log.Fatal(err)
	}
	err = bot.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
