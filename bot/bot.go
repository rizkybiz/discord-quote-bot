package bot

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rizkybiz/discord-quote-bot/commands"
	"github.com/rizkybiz/discord-quote-bot/quotes"
)

//Bot is a struct that provides top level control for the discord bot
type Bot struct {
	session  *discordgo.Session
	commands commands.Commands
	quotes   quotes.Quotes
}

//New returns a new bot with discord session
// created as well as the commands and quotes structures
func New(token string) (Bot, error) {

	// Create the discord session with the bot token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return Bot{}, err
	}
	// Set up the commands and quotes
	c := commands.New()
	q := quotes.New()
	commands.RegisterHelper(c, q)
	dg.AddHandler(c.ProcessCmd)

	return Bot{
		session:  dg,
		commands: c,
		quotes:   q,
	}, nil
}

//Connect establishes a connection for the bot to the Discord API
func (b Bot) Connect() error {
	// Connect to the server as the bot
	err := b.session.Open()
	if err != nil {
		return err
	}
	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Term signal has been received, cleanup and close.
	b.session.Close()
	return nil
}
