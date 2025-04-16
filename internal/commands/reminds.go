package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Constants for command parameters
const (
	whenParam = "--when"
	textParam = "--text"
	idParam   = "--id"
)

// NewAddRemindCommand creates a command to add a reminder
func NewAddRemindCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-add",
		Description: "Add a new reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The reminder text",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Adding reminder")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No reminder text provided.")
			return
		}

		text := options[0].StringValue()
		if strings.TrimSpace(text) == "" {
			respondError(s, i, "Reminder text cannot be empty.")
			return
		}

		requestBody := buildBaseRequestBody(i, "/todo-remind")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		requestBody["text"] = text

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionAdd, "Failed to add your reminder: ")
	}

	return cmd, handler
}

// NewAddOneRemindCommand creates a command to add a reminder with timing
func NewAddOneRemindCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-add-one",
		Description: "Add a new reminder with specific timing",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "when",
				Description: "Cron timing setting (e.g., '0 9 * * 1' for every Monday at 9am)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Reminder text",
				Required:    true, // Changed to true for clarity
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Adding timed reminder")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options
		if len(options) < 2 {
			respondError(s, i, "Both timing and reminder text are required.")
			return
		}

		when := strings.TrimSpace(options[0].StringValue())
		if when == "" {
			respondError(s, i, "Timing parameter cannot be empty.")
			return
		}

		text := strings.TrimSpace(options[1].StringValue())
		if text == "" {
			respondError(s, i, "Reminder text cannot be empty.")
			return
		}

		// Properly format the command parameters with spaces
		todoText := fmt.Sprintf("%s=%s %s=%s", whenParam, when, textParam, text)

		requestBody := buildBaseRequestBody(i, "/todo-remind-add-one")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		requestBody["text"] = todoText

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionAdd, "Failed to add your timed reminder: ")
	}

	return cmd, handler
}

// NewListRemindsCommand creates a command to list all reminders
func NewListRemindsCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-list",
		Description: "List all your reminders",
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Listing reminders")

		requestBody := buildBaseRequestBody(i, "/todo-remind-list")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionList, "Failed to list your reminders: ")
	}

	return cmd, handler
}

func NewEnableRemindsCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-enable",
		Description: "Enable reminders",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the reminder",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		options := i.ApplicationCommandData().Options
		id := strings.TrimSpace(options[0].StringValue())

		todoText := fmt.Sprintf("%s=%s", idParam, id)

		requestBody := buildBaseRequestBody(i, "/todo-remind-enable")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		requestBody["text"] = todoText

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionEnable, "Failed to enable your reminder: ")
	}

	return cmd, handler
}

func NewDisableRemindsCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-disable",
		Description: "Disable reminders",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the reminder",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		id := strings.TrimSpace(options[0].StringValue())

		todoText := fmt.Sprintf("%s=%s", idParam, id)

		requestBody := buildBaseRequestBody(i, "/todo-remind-disable")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		requestBody["text"] = todoText

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionDisable, "Failed to disable your reminder: ")
	}

	return cmd, handler
}

func NewUpdateRemindCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-update",
		Description: "Update reminders",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "The ID of the reminder",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "when",
				Description: "Cron timing setting (e.g., '0 9 * * 1' for every Monday at 9am)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "Reminder text",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		id := strings.TrimSpace(options[0].StringValue())

		todoText := fmt.Sprintf("%s=%s", idParam, id)

		requestBody := buildBaseRequestBody(i, "/todo-remind-update")
		if requestBody == nil {
			respondError(s, i, "Failed to build request data.")
			return
		}

		requestBody["text"] = todoText

		endpoint := fmt.Sprintf("%s/reminds", apiBaseURL)
		executeTodoRequest(s, i, endpoint, requestBody, RemindActionUpdate, "Failed to update your reminder: ")
	}

	return cmd, handler
}
