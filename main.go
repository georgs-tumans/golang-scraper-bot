package main

import (
	"fmt"
	"log"
	services "web_scraper_bot/services/bot"
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
	botService := services.NewBotService()
	botService.InitializeBot()
}
