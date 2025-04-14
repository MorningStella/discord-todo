package commands

import (
	"log"

	"github.com/MorningStella/discord-todo/internal/api"
	"github.com/bwmarrin/discordgo"
)

// CommandHandler is a function type that handles Discord interactions

// buildBaseRequestBody creates a common request structure from interaction data
func buildBaseRequestBody(
	i *discordgo.InteractionCreate,
	commandName string,
) map[string]interface{} {
	return map[string]interface{}{
		"channel_id":   i.ChannelID,
		"channel_name": "", // Fill this if available from interaction
		"command":      commandName,
		"response_url": "", // Fill this if available from interaction
		"team_domain":  "", // Fill this if available from interaction
		"team_id":      i.GuildID,
		"token":        "", // Fill this if available from interaction
		"trigger_id":   i.ID,
		"user_id":      i.Member.User.ID,
		"user_name":    i.Member.User.Username,
	}
}

// executeTodoRequest executes a todo API request and handles response/errors
func executeTodoRequest(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	apiBaseURL string,
	requestBody map[string]interface{},
	todoAction TodoAction,
	errorPrefix string,
) {
	apiClient := api.NewClient(apiBaseURL)
	message, err := apiClient.SendTodoRequest(requestBody, todoAction.String())
	if err != nil {
		log.Printf("Error with %s operation: %v", todoAction, err)
		respondError(s, i, errorPrefix+err.Error())
		return
	}

	respond(s, i, message)
}

// NewAddTodoCommand creates a new todo command definition and handler
func NewAddTodoCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "todo-add",
		Description: "Add a new todo item",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The todo item text",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Adding todo...")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No todo item text provided.")
			return
		}

		todoText := options[0].StringValue()
		requestBody := buildBaseRequestBody(i, "/todo-add")
		requestBody["text"] = todoText

		executeTodoRequest(s, i, apiBaseURL, requestBody, TodoActionAdd, "Failed to add your todo item: ")
	}

	return cmd, handler
}

// NewListTodoCommand creates a command to list todos
func NewListTodoCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "todo-list",
		Description: "List current channel todo items",
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Listing todo...")

		requestBody := buildBaseRequestBody(i, "/todo-list")
		executeTodoRequest(s, i, apiBaseURL, requestBody, TodoActionList, "Failed to list your todo items: ")
	}

	return cmd, handler
}

// NewCompleteTodoCommand creates a command to complete todos
func NewCompleteTodoCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "todo-done",
		Description: "Complete todo by id",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The todo item ID to mark as complete",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Completing todo...")

		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No todo item ID provided.")
			return
		}

		todoID := options[0].StringValue()
		requestBody := buildBaseRequestBody(i, "/todo-done")
		requestBody["text"] = todoID

		executeTodoRequest(s, i, apiBaseURL, requestBody, TodoActionDone, "Failed to complete your todo item: ")
	}

	return cmd, handler
}

func NewUpdateTodoCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "todo-update",
		Description: "Update todo by id",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Update todo command",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Completing todo...")

		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No todo item ID provided.")
			return
		}

		text := options[0].StringValue()
		requestBody := buildBaseRequestBody(i, "/todo-update")
		requestBody["text"] = text

		executeTodoRequest(s, i, apiBaseURL, requestBody, TodoActionUpdate, "Failed to complete your todo item: ")
	}

	return cmd, handler
}

// respond sends a response to a Discord interaction
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
	}
}

// respondError sends an error response to a Discord interaction
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + message,
		},
	})

	if err != nil {
		log.Printf("Error sending error response: %v", err)
	}
}
