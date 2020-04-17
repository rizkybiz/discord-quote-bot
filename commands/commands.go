package commands

import "github.com/bwmarrin/discordgo"

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
