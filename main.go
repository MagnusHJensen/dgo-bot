package main

import (
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	"magnusjensen.dk/dgo-bot/internal"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		internal.LOGGER.Fatal("DISCORD_BOT_TOKEN environment variable not set")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		internal.LOGGER.Fatal("Error creating Discord session:", err)
	}

	err = discord.Open()
	if err != nil {
		internal.LOGGER.Fatal("Error opening Discord session:", err)
	}

	// Bot is now ready.
	internal.LOGGER.Println("Bot is now running.")

	// Event Handlers
	discord.AddHandler(internal.HandleMCLogsMessage)

	// Wait here until CTRL-C or other term signal is received.
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = discord.Close()
	if err != nil {
		internal.LOGGER.Printf("could not close session gracefully: %s", err)
	}
}
