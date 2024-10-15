package main

import (
	"fmt"
	"log"
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
	//botService.InitializeBotLongPolling()
	botFixer.InitializeBotWebhook()
}
