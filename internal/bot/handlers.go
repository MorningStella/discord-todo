package bot

import (
	"github.com/MorningStella/discord-todo/internal/commands"
)

// registerCommands registers all command definitions and handlers
func (b *Bot) registerCommands() {
	// Get todo command
	todoCmd, todoHandler := commands.NewAddTodoCommand(b.config.APIBaseURL)
	listTodoCmd, listTodoHandler := commands.NewListTodoCommand(b.config.APIBaseURL)

	// Add commands to the bot's command list
	b.commands = append(b.commands, todoCmd)
	b.commands = append(b.commands, listTodoCmd)

	// Register command handlers
	b.commandHandlers[todoCmd.Name] = todoHandler
	b.commandHandlers[listTodoCmd.Name] = listTodoHandler

	// Add more commands here as needed:
	// someCmd, someHandler := commands.NewSomeCommand()
	// b.commands = append(b.commands, someCmd)
	// b.commandHandlers[someCmd.Name] = someHandler
}
