package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MorningStella/discord-todo/internal/bot"
	"github.com/MorningStella/discord-todo/internal/config"
)

func main() {
	// Parse command line flags
	cfg := config.LoadConfig()
	fmt.Println(cfg.BotToken)
	// Initialize the bot
	discordBot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Start the bot
	if err := discordBot.Start(); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
	defer discordBot.Close()

	// Wait for interrupt signal to gracefully shutdown
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Cleanup before exit
	if cfg.RemoveCommands {
		discordBot.RemoveCommands()
	}

	log.Println("Gracefully shutting down.")
}
