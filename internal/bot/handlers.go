package bot

import (
	"github.com/MorningStella/discord-todo/internal/api"
	"github.com/MorningStella/discord-todo/internal/commands"
)

// registerCommands registers all command definitions and handlers
func (b *Bot) registerCommands(workflowClient *api.WorkflowClient, reclaimClient *api.ReclaimClient) {
	// Get todo command
	todoCmd, todoHandler := commands.NewAddTodoCommand(workflowClient)
	listTodoCmd, listTodoHandler := commands.NewListTodoCommand(workflowClient)
	completeTodoCmd, completeTodoHandler := commands.NewCompleteTodoCommand(workflowClient)
	updateTodoCmd, updateTodoHandler := commands.NewUpdateTodoCommand(workflowClient)

	addOneRemindCmd, addOneRemindHandler := commands.NewAddOneRemindCommand(workflowClient)
	listRemindCmd, listRemindHandler := commands.NewListRemindsCommand(workflowClient)
	// Enable and disable remind commands
	enableRemindCmd, enableRemindHandler := commands.NewEnableRemindsCommand(workflowClient)
	disableRemindCmd, disableRemindHandler := commands.NewDisableRemindsCommand(workflowClient)

	updateRemindCmd, updateRmindaddOneRemindHandler := commands.NewUpdateRemindCommand(workflowClient)

	// Reclaim task
	reclaimTaskCmd, reclaimTaskHandler := commands.NewCreateReclaimTaskCommand(reclaimClient)

	// Add commands to the bot's command list
	b.commands = append(b.commands, todoCmd)
	b.commands = append(b.commands, listTodoCmd)
	b.commands = append(b.commands, completeTodoCmd)
	b.commands = append(b.commands, updateTodoCmd)
	b.commands = append(b.commands, addOneRemindCmd)
	b.commands = append(b.commands, listRemindCmd)
	b.commands = append(b.commands, enableRemindCmd)
	b.commands = append(b.commands, disableRemindCmd)
	b.commands = append(b.commands, updateRemindCmd)
	b.commands = append(b.commands, reclaimTaskCmd)

	// Register command handlers
	b.commandHandlers[todoCmd.Name] = todoHandler
	b.commandHandlers[listTodoCmd.Name] = listTodoHandler
	b.commandHandlers[completeTodoCmd.Name] = completeTodoHandler
	b.commandHandlers[updateTodoCmd.Name] = updateTodoHandler
	b.commandHandlers[addOneRemindCmd.Name] = addOneRemindHandler
	b.commandHandlers[listRemindCmd.Name] = listRemindHandler
	b.commandHandlers[enableRemindCmd.Name] = enableRemindHandler
	b.commandHandlers[disableRemindCmd.Name] = disableRemindHandler
	b.commandHandlers[updateRemindCmd.Name] = updateRmindaddOneRemindHandler
	b.commandHandlers[reclaimTaskCmd.Name] = reclaimTaskHandler
}
