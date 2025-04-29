package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MorningStella/discord-todo/internal/api"

	"github.com/bwmarrin/discordgo"
)

// NewCreateReclaimTaskCommand creates a command to add a task to Reclaim.ai
func NewCreateReclaimTaskCommand(reclaimClient *api.ReclaimClient) (*discordgo.ApplicationCommand, CommandHandler) {
	cmd := &discordgo.ApplicationCommand{
		Name:        "task-create",
		Description: "Create a new task in Reclaim.ai",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "Title of the task",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "category",
				Description: "Category: 'work' or 'personal'",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Work",
						Value: "WORK",
					},
					{
						Name:  "Personal",
						Value: "PERSONAL",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "priority",
				Description: "Priority level (1-4, where 1 is highest)",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "P1 - Highest",
						Value: 1,
					},
					{
						Name:  "P2 - High",
						Value: 2,
					},
					{
						Name:  "P3 - Medium",
						Value: 3,
					},
					{
						Name:  "P4 - Low",
						Value: 4,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "due-date",
				Description: "Due date (format: YYYY-MM-DD)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "time-required",
				Description: "Time required in minutes",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "min-chunk",
				Description: "Minimum time chunk size in minutes",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "max-chunk",
				Description: "Maximum time chunk size in minutes",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "private",
				Description: "Make task always private",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "on-deck",
				Description: "Add task to 'On Deck'",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "schedule-after",
				Description: "Don't schedule before this date (format: YYYY-MM-DD)",
				Required:    false,
			},
		},
	}

	// Command handler
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Printf("Creating Reclaim task")

		// Get user input from the command options
		options := i.ApplicationCommandData().Options

		// Default values for optional parameters
		var (
			title            string
			categoryStr      string
			priorityNum      int64
			dueDateStr       string
			timeRequired     int64
			minChunk         int64 = 30    // Default: 30 minutes
			maxChunk         int64 = 120   // Default: 2 hours
			private          bool  = false // Default: not private
			onDeck           bool  = false // Default: not on deck
			scheduleAfterStr string
		)

		// Extract all parameters from options
		for _, opt := range options {
			switch opt.Name {
			case "title":
				title = strings.TrimSpace(opt.StringValue())
			case "category":
				categoryStr = opt.StringValue()
			case "priority":
				priorityNum = opt.IntValue()
			case "due-date":
				dueDateStr = strings.TrimSpace(opt.StringValue())
			case "time-required":
				timeRequired = opt.IntValue()
			case "min-chunk":
				minChunk = opt.IntValue()
			case "max-chunk":
				maxChunk = opt.IntValue()
			case "private":
				private = opt.BoolValue()
			case "on-deck":
				onDeck = opt.BoolValue()
			case "schedule-after":
				scheduleAfterStr = strings.TrimSpace(opt.StringValue())
			}
		}

		// Parameter validation
		if title == "" {
			respondError(s, i, "Task title cannot be empty.")
			return
		}

		// Validate time-required fits in uint8
		if timeRequired > 255 {
			respondError(s, i, "Time required cannot exceed 255 minutes.")
			return
		}

		// Validate chunk sizes fit in uint8
		if minChunk > 255 || maxChunk > 255 {
			respondError(s, i, "Chunk sizes cannot exceed 255 minutes.")
			return
		}

		// Parse category
		var category api.EventCategory
		if categoryStr == "WORK" {
			category = api.Work
		} else if categoryStr == "PERSONAL" {
			category = api.Personal
		} else {
			respondError(s, i, "Invalid category. Use 'work' or 'personal'.")
			return
		}

		// Map priority number to TaskPriority
		var priority api.TaskPriority
		switch priorityNum {
		case 1:
			priority = api.PriorityP1
		case 2:
			priority = api.PriorityP2
		case 3:
			priority = api.PriorityP3
		case 4:
			priority = api.PriorityP4
		default:
			respondError(s, i, "Invalid priority. Use a number from 1 to 4.")
			return
		}

		// Parse due date
		dueDate, err := time.Parse("2006-01-02", dueDateStr)
		if err != nil {
			respondError(s, i, "Invalid due date format. Use YYYY-MM-DD.")
			return
		}

		// Set due date to end of day (now a pointer in the struct)
		dueDate = time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 23, 59, 59, 0, time.Local)

		// Parse schedule after date if provided
		var scheduleAfter *time.Time
		if scheduleAfterStr != "" {
			parsedDate, err := time.Parse("2006-01-02", scheduleAfterStr)
			if err != nil {
				respondError(s, i, "Invalid schedule-after date format. Use YYYY-MM-DD.")
				return
			}
			scheduleAfter = &parsedDate
		}

		// Create the task
		task := api.ReclaimTask{
			Title:              &title,
			EventCategory:      category,
			TimeChunksRequired: uint8(timeRequired),
			AlwaysPrivate:      private,
			Priority:           priority,
			ScheduleAfter:      scheduleAfter,
			Due:                &dueDate,
			OnDeck:             onDeck,
		}

		// Only set minChunk and maxChunk if they were provided
		if minChunk != 30 { // If not default value
			minChunkUint8 := uint8(minChunk)
			task.MinChunkSize = &minChunkUint8
		}

		if maxChunk != 120 { // If not default value
			maxChunkUint8 := uint8(maxChunk)
			task.MaxChunkSize = &maxChunkUint8
		}

		// Call the Reclaim API
		response, err := reclaimClient.CreateTask(task)
		if err != nil {
			respondError(s, i, fmt.Sprintf("Failed to create Reclaim task: %v", err))
			return
		}

		fmt.Println("successfully created task:", response)

		// Prepare success message
		successMsg := fmt.Sprintf("✅ Created task: **%s**\n", title)
		successMsg += fmt.Sprintf("• Category: %s\n", categoryStr)
		successMsg += fmt.Sprintf("• Priority: P%d\n", priorityNum)
		successMsg += fmt.Sprintf("• Due: %s\n", dueDateStr)

		// Add task ID or URL if available in response
		if taskID, ok := response["id"].(string); ok {
			successMsg += fmt.Sprintf("• Task ID: %s\n", taskID)
		}

		// Send the response
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &successMsg,
		})
	}

	return cmd, handler
}
