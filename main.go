package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rizkybiz/discord-quote-bot/commands"
)

//Token is the token of the dicord bot
var Token string

//Setup the commands map globally
var commandsMap = commands.New()

func init() {
	flag.StringVar(&Token, "token", "", "token of the discord bot")
	flag.Parse()
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

	//Register the commandsMap triggers and handler functions
	commandsMap.Register("!help", func(s *discordgo.Session, m *discordgo.Message) {
		helpText := `List of Available Commands:
		!help - Displays this message
		!addquote - Use this command followed by a space, @<user> another space, then the quote you'd like to add
		!quote - Use this command followed by a space and @<user> for a random quote from your favorite user`
		s.ChannelMessageSend(m.ChannelID, helpText)
	})

	// Register the commands handler function from below so the session knows what to do
	dg.AddHandler(processCmd)

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

func processCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Check if this is even a trigger command
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	//Extract the trigger command
	msgList := strings.SplitAfter(m.Content, " ")
	trigger := msgList[0]

	//Get the message from the messageCreate event
	message := m.Message

	//Check for existence of the trigger command and corresponding handler func
	//If it doesn't exist, return the help message
	if cmdFunc, ok := commandsMap[trigger]; ok {
		cmdFunc(s, message)
	} else {
		log.Printf("Command %s does not exist", trigger)
		s.ChannelMessageSend(m.ChannelID, `That command does not exist, type "!help" for available commands`)
	}
}
