package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rizkybiz/discord-quote-bot/quotes"
)

//Command is a function signature with all the info we need to run commands
type Command func(session *discordgo.Session, message *discordgo.Message)

//Commands is a map to store triggers mapped to handler functions
type Commands map[string]Command

//New creates a new map of triggers to handler functions
func New() Commands {
	return make(map[string]Command)
}

//Register adds a new trigger and handler function to the Commands map
func (c Commands) Register(name string, handler Command) {
	c[name] = handler
}

//Unregister removes a trigger and handler function from the Commands map
func (c Commands) Unregister(name string) {
	delete(c, name)
}

//RegisterHelper registers the ! commands of the bot
func RegisterHelper(c Commands, q quotes.Quotes) {
	//Register the commandsMap triggers and handler functions
	c.Register("!help", func(s *discordgo.Session, m *discordgo.Message) {
		helpText := `List of Available Commands:
		!help - Displays this message
		!addquote - Use this command followed by a space, @<user> another space, then the quote you'd like to add
		!quote - Use this command followed by a space and @<user> for a random quote from your favorite user`
		s.ChannelMessageSend(m.ChannelID, helpText)
	})

	c.Register("!addquote", func(s *discordgo.Session, m *discordgo.Message) {
		if len(m.Mentions) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No user mention included. Please @<user> type a space then type a quote")
			log.Println("No user mentioned, couldn't add quote")
			return
		}
		user := m.Mentions[0].ID
		parsedMessage := strings.SplitAfter(m.Content, " ")
		if len(parsedMessage) < 3 {
			s.ChannelMessageSend(m.ChannelID, "No quote mention included. Please @<user> type a space then type a quote")
			log.Println("No quote attached, couldn't add quote")
			return
		}
		var quote string
		for i := 2; i <= len(parsedMessage)-1; i++ {
			quote = quote + parsedMessage[i]
		}
		q.Add(user, quote)
		s.ChannelMessageSend(m.ChannelID, "Quote added!")
		log.Printf("New quote added for %s", user)
		return
	})

	c.Register("!quote", func(s *discordgo.Session, m *discordgo.Message) {
		if len(m.Mentions) == 0 {
			s.ChannelMessageSend(m.ChannelID, "No user mention included. Please @<user>")
			log.Println("No user mentioned, couldn't get quote")
			return
		}
		user := m.Mentions[0]
		quote, err := q.Get(user.ID)
		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(`"%s" - <@%s>`, quote, user.ID))
	})
}

//ProcessCmd is a handler function that serves as the entrypoint
// for ! commands
func (c Commands) ProcessCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	triggerStr := strings.TrimSpace(trigger)

	//Get the message from the messageCreate event
	message := m.Message

	//Check for existence of the trigger command and corresponding handler func
	//If it doesn't exist, return the help message
	if cmdFunc, ok := c[triggerStr]; ok {
		cmdFunc(s, message)
	} else {
		log.Printf("Command %s does not exist", triggerStr)
		s.ChannelMessageSend(m.ChannelID, `That command does not exist, type "!help" for available commands`)
	}
}
