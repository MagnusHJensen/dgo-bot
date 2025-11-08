package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN environment variable not set")
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	fmt.Println("Bot is now running.")

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	err = discord.Close()
	if err != nil {
		log.Printf("could not close session gracefully: %s", err)
	}
}
