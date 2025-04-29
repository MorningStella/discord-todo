package bot

import (
	"log"

	"github.com/MorningStella/discord-todo/internal/api"
	"github.com/MorningStella/discord-todo/internal/config"
	"github.com/bwmarrin/discordgo"
)

// Bot represents the Discord bot
type Bot struct {
	session         *discordgo.Session
	config          *config.Config
	commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	commands        []*discordgo.ApplicationCommand
	registeredCmds  []*discordgo.ApplicationCommand
}

// NewDiscordBot creates a new bot instance
func NewDiscordBot(cfg *config.Config, workflowClient *api.WorkflowClient, reclaimClient *api.ReclaimClient) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		session:         session,
		config:          cfg,
		commandHandlers: make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
	}

	// Register commands
	bot.registerCommands(workflowClient, reclaimClient)

	return bot, nil
}

// Start starts the Discord bot
func (b *Bot) Start() error {
	// Add ready handler
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Add interaction handler
	b.session.AddHandler(b.handleInteraction)

	// Open connection to Discord
	if err := b.session.Open(); err != nil {
		return err
	}

	// Register slash commands
	b.removeExistingCommands()
	b.registerSlashCommands()

	return nil
}

// Close closes the Discord session
func (b *Bot) Close() {
	b.session.Close()
}

// handleInteraction handles Discord interactions
func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Only handle ApplicationCommand interactions
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	// Find and execute the appropriate command handler
	if handler, ok := b.commandHandlers[i.ApplicationCommandData().Name]; ok {
		handler(s, i)
	}
}

// removeExistingCommands removes existing slash commands
func (b *Bot) removeExistingCommands() {
	log.Println("Removing existing commands...")

	existingCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, b.config.GuildID)
	if err != nil {
		log.Printf("Could not fetch registered commands: %v", err)
		return
	}

	for _, cmd := range existingCommands {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.config.GuildID, cmd.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}

// registerSlashCommands registers all slash commands
func (b *Bot) registerSlashCommands() {
	log.Println("Registering commands...")

	b.registeredCmds = make([]*discordgo.ApplicationCommand, len(b.commands))

	for i, cmd := range b.commands {
		registeredCmd, err := b.session.ApplicationCommandCreate(
			b.session.State.User.ID,
			b.config.GuildID,
			cmd,
		)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", cmd.Name, err)
		}
		b.registeredCmds[i] = registeredCmd
	}
}

// RemoveCommands removes all registered commands
func (b *Bot) RemoveCommands() {
	log.Println("Removing commands...")

	for _, cmd := range b.registeredCmds {
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, b.config.GuildID, cmd.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
