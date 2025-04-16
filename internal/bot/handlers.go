package bot

import (
	"github.com/MorningStella/discord-todo/internal/commands"
)

// registerCommands registers all command definitions and handlers
func (b *Bot) registerCommands() {
	// Get todo command
	todoCmd, todoHandler := commands.NewAddTodoCommand(b.config.APIBaseURL)
	listTodoCmd, listTodoHandler := commands.NewListTodoCommand(b.config.APIBaseURL)
	completeTodoCmd, completeTodoHandler := commands.NewCompleteTodoCommand(b.config.APIBaseURL)
	updateTodoCmd, updateTodoHandler := commands.NewUpdateTodoCommand(b.config.APIBaseURL)

	addRemindCmd, addRemindHandler := commands.NewAddRemindCommand(b.config.APIBaseURL)
	addOneRemindCmd, addOneRemindHandler := commands.NewAddOneRemindCommand(b.config.APIBaseURL)
	listRemindCmd, listRemindHandler := commands.NewListRemindsCommand(b.config.APIBaseURL)

	// Add commands to the bot's command list
	b.commands = append(b.commands, todoCmd)
	b.commands = append(b.commands, listTodoCmd)
	b.commands = append(b.commands, completeTodoCmd)
	b.commands = append(b.commands, updateTodoCmd)
	b.commands = append(b.commands, addRemindCmd)
	b.commands = append(b.commands, addRemindCmd)
	b.commands = append(b.commands, addOneRemindCmd)
	b.commands = append(b.commands, listRemindCmd)

	// Register command handlers
	b.commandHandlers[todoCmd.Name] = todoHandler
	b.commandHandlers[listTodoCmd.Name] = listTodoHandler
	b.commandHandlers[completeTodoCmd.Name] = completeTodoHandler
	b.commandHandlers[updateTodoCmd.Name] = updateTodoHandler
	b.commandHandlers[addRemindCmd.Name] = addRemindHandler
	b.commandHandlers[addOneRemindCmd.Name] = addOneRemindHandler
	b.commandHandlers[listRemindCmd.Name] = listRemindHandler

	// Add more commands here as needed:
	// someCmd, someHandler := commands.NewSomeCommand()
	// b.commands = append(b.commands, someCmd)
	// b.commandHandlers[someCmd.Name] = someHandler
}
