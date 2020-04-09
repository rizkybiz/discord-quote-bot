package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

//Token is the token of the dicord bot
var Token string

func init() {
	flag.StringVar(&Token, "token", "", "token of the discord bot")
}

func main() {

	// Check if the token has been provided
	if Token == "" {
		log.Fatalf("Discord bot token not set, exiting")
	}

	// Create the discord session with the bot token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln(err)
	}

	// Register the pingPongHandler function from below so the session knows what to do
	dg.AddHandler(pingPongHandler)

	// Connect to the server as the bot
	err = dg.Open()
	if err != nil {
		log.Fatalln(err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Term signal has been received, cleanup and close.
	dg.Close()
	return
}

//pingPongHandler is a handler function which will be registered in the main function on the *discordgo.Session
func pingPongHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
