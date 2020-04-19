package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rizkybiz/discord-quote-bot/commands"
	"github.com/rizkybiz/discord-quote-bot/quotes"
)

//Token is the token of the dicord bot
var Token string

//Setup the commands map globally
var commandsMap = commands.New()

//Setup the map of users to quotes
var quotesMap = quotes.New()

func init() {
	flag.StringVar(&Token, "token", "", "token of the discord bot")
	flag.Parse()
}

func main() {
	//TODO: Move the handler registration out of here into the commands package for readability, also processCmd

	// Check if the token has been provided
	if Token == "" {
		log.Fatalf("Discord bot token not set, exiting")
	}

	// Create the discord session with the bot token
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln(err)
	}

	//Register the ! commands onto the commandsMap
	commands.RegisterHelper(commandsMap, quotesMap)

	// Register the commands handler function so the session knows what to do
	dg.AddHandler(commandsMap.ProcessCmd)

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
