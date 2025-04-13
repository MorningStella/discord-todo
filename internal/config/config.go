package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	GuildID        string
	BotToken       string
	RemoveCommands bool
	APIBaseURL     string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it. Using environment variables directly.")
	}

	cfg := &Config{}

	// Read from environment variables
	cfg.GuildID = os.Getenv("GUILD_ID")
	cfg.BotToken = os.Getenv("BOT_TOKEN")

	// Default to true for removing commands
	cfg.RemoveCommands = true

	// Override if environment variable is provided
	if removeStr := os.Getenv("REMOVE_COMMANDS"); removeStr != "" {
		removeBool, err := strconv.ParseBool(removeStr)
		if err == nil {
			cfg.RemoveCommands = removeBool
		}
	}

	// Set API base URL with a fallback
	cfg.APIBaseURL = os.Getenv("API_BASE_URL")
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = "https://n8n.weiting.me/webhook-test/discord"
	}

	// Validate required configuration
	if cfg.BotToken == "" {
		log.Println("Warning: BOT_TOKEN environment variable is not set")
	}

	return cfg
}
