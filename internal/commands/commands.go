package commands

import (
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function that handles a Discord interaction
type CommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Command represents a Discord slash command with its handler
type Command struct {
	Definition *discordgo.ApplicationCommand
	Handler    CommandHandler
}
