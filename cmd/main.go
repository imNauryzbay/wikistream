package main

import (
	"flag"
	"log"

	"wikistream/internal/bot"
)

func main() {
	token := flag.String("token", "", "Bot authentication token")
	flag.Parse()

	if *token == "" {
		log.Fatal("Bot token is not set")
	}

	discordBot, err := bot.New(*token)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	if err := discordBot.Run(); err != nil {
		log.Fatalf("Bot stopped with error: %v", err)
	}
}
