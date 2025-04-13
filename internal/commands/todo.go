package commands

import (
	"log"

	"github.com/MorningStella/discord-todo/internal/api"
	"github.com/bwmarrin/discordgo"
)

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
		username := i.Member.User.Username
		userId := i.Member.User.ID
		channelId := i.ChannelID

		// Build the request based on command data
		requestBody := map[string]interface{}{
			"channel_id":   channelId,
			"channel_name": "", // Fill this if available from interaction
			"command":      "/todo-add",
			"response_url": "", // Fill this if available from interaction
			"team_domain":  "", // Fill this if available from interaction
			"team_id":      i.GuildID,
			"text":         todoText,
			"token":        "", // Fill this if available from interaction
			"trigger_id":   i.ID,
			"user_id":      userId,
			"user_name":    username,
		}

		apiClient := api.NewClient(apiBaseURL)
		message, err := apiClient.SendTodoRequest(requestBody, "add")
		if err != nil {
			log.Printf("Error adding todo: %v", err)
			respondError(s, i, "Failed to add your todo item: "+err.Error())
			return
		}

		// Send success response with the actual message from the API
		respond(s, i, message)
	}

	return cmd, handler
}

func NewListTodoCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "todo-list",
		Description: "List current channel todo items",
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Listing todo...")

		username := i.Member.User.Username
		userId := i.Member.User.ID
		channelId := i.ChannelID

		// Build the request based on command data
		requestBody := map[string]interface{}{
			"channel_id":   channelId,
			"channel_name": "", // Fill this if available from interaction
			"command":      "/todo-list",
			"response_url": "", // Fill this if available from interaction
			"team_domain":  "", // Fill this if available from interaction
			"team_id":      i.GuildID,
			"token":        "", // Fill this if available from interaction
			"trigger_id":   i.ID,
			"user_id":      userId,
			"user_name":    username,
		}

		apiClient := api.NewClient(apiBaseURL)
		message, err := apiClient.SendTodoRequest(requestBody, "list")
		if err != nil {
			log.Printf("Error adding todo: %v", err)
			respondError(s, i, "Failed to add your todo item: "+err.Error())
			return
		}

		// Send success response with the actual message from the API
		respond(s, i, message)
	}

	return cmd, handler
}

// respond sends a response to a Discord interaction
func respond(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

// respondError sends an error response to a Discord interaction
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Error: " + message,
		},
	})
}
