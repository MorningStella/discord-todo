package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func NewAddRemindCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-add",
		Description: "Add a new reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The remind setting",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Adding remind")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No remind item text provided.")
			return
		}

		todoText := options[0].StringValue()
		requestBody := buildBaseRequestBody(i, "/todo-remind")
		requestBody["text"] = todoText

		executeTodoRequest(s, i, fmt.Sprintf("%s/reminds", apiBaseURL), requestBody, RemindActionAdd, "Failed to add your todo item: ")
	}

	return cmd, handler
}

func NewAddOneRemindCommand(apiBaseURL string) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "remind-add-one",
		Description: "Add a new reminder",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The remind setting",
				Required:    true,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Adding remind")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options
		if len(options) == 0 {
			respondError(s, i, "No remind item text provided.")
			return
		}

		todoText := options[0].StringValue()
		requestBody := buildBaseRequestBody(i, "/todo-remind-add-one")
		requestBody["text"] = todoText

		executeTodoRequest(s, i, fmt.Sprintf("%s/reminds", apiBaseURL), requestBody, RemindActionAdd, "Failed to add your todo item: ")
	}

	return cmd, handler
}
