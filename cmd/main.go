package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/MorningStella/discord-todo/internal/api"
	"github.com/MorningStella/discord-todo/internal/bot"
	"github.com/MorningStella/discord-todo/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println(cfg.BotToken)
	workflowClient := api.NewWorkflowClient(cfg)
	reclaimClient := api.NewReclaimClient(cfg)
	discordBot, err := bot.NewDiscordBot(cfg, workflowClient, reclaimClient)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	if err := discordBot.Start(); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}
	defer discordBot.Close()

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
