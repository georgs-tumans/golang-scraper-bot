package bot_fixer

import (
	"context"
	"log"
	"net/http"
	"web_scraper_bot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotFixer struct {
	Bot               *tgbotapi.BotAPI
	Config            *config.Configuration
	BondsClientActive bool
}

func NewBotFixer() *BotFixer {
	botService := &BotFixer{
		Config: config.GetConfig(),
	}

	var err error
	botService.Bot, err = tgbotapi.NewBotAPI(botService.Config.BotAPIKey)
	if err != nil {
		log.Panic(err)
		return nil
	}

	return botService
}

func (b *BotFixer) InitializeBotLongPolling() {
	// Set this to true to log all interactions with telegram servers
	b.Bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()

	// `updates` is a golang channel which receives telegram updates
	updates := b.Bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go b.receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("[Bot fixer] Listening for updates.")

	select {}
}

func (b *BotFixer) InitializeBotWebhook() {
	// Set this to true to log all interactions with telegram servers
	b.Bot.Debug = false

	wh, err := tgbotapi.NewWebhook(b.Config.WebhookURL + "/webhook")
	if err != nil {
		log.Fatalf("[Bot fixer] Error creating webhook: %v", err)
	}

	_, err = b.Bot.Request(wh)
	if err != nil {
		log.Fatalf("[Bot fixer] Error setting webhook: %v", err)
	}

	log.Printf("[Bot fixer] Webhook set: %s", b.Config.WebhookURL+"/webhook")

	http.HandleFunc("/webhook", b.webhookHandler)

	log.Println("[Bot fixer] Starting server on port " + b.Config.Port)
	log.Fatal(http.ListenAndServe(":"+b.Config.Port, nil))
}
