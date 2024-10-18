package main

import (
	"fmt"
	"log"
	"strings"
	"web_scraper_bot/bot_fixer"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("unhandled error: %v", r)
			log.Fatalf("[main] Unhandled panic: %v", err)

			return
		}
	}()

	log.Println("Starting bot service")
	botFixer := bot_fixer.NewBotFixer()
	config := botFixer.Config

	if strings.ToLower(strings.TrimSpace(config.Environment)) == "local" {
		if err := botFixer.DeleteWebhook(); err != nil {
			log.Fatalf("[main] Error deleting webhook: %v", err)
		}
		botFixer.InitializeBotLongPolling()
	} else {
		botFixer.InitializeBotWebhook()
	}
}
