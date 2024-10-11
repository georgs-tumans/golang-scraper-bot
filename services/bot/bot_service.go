package bot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"web_scraper_bot/clients"
	"web_scraper_bot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	Bot               *tgbotapi.BotAPI
	Config            *config.Configuration
	BondsClientActive bool
}

func NewBotService() *BotService {
	botService := &BotService{
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

func (b *BotService) InitializeBot() {
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
	log.Println("[Bot service] Listening for updates.")

	select {}
}

func (b *BotService) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			b.handleUpdate(update)
		}
	}
}

func (b *BotService) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		b.handleMessage(update.Message)

	// Handle button clicks
	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery)
	}
}

func (b *BotService) handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("[Bot service] %s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = b.handleCommand(message.Chat.ID, text)
	}

	if err != nil {
		log.Printf("[Bot service] An error occured while handlind message: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func (b *BotService) handleCommand(chatId int64, command string) error {
	var err error

	switch command {

	case "/menu":
		err = b.SendMenu(chatId)

	case "/bonds_start":
		if !b.BondsClientActive {
			b.BondsClientActive = true
			b.SendMessage(chatId, "Savings bonds client has been started", nil)
			go b.activateBondsClient(chatId) // Run in a separate goroutine
		} else {
			b.SendMessage(chatId, "Savings bonds client is already running", nil)
		}

	case "/bonds_stop":
		if b.BondsClientActive {
			b.BondsClientActive = false
			b.SendMessage(chatId, "Savings bonds client has been stopped", nil)
		} else {
			b.SendMessage(chatId, "Savings bonds client is not running", nil)
		}
	}

	return err
}

func (b *BotService) activateBondsClient(chatId int64) {
	ticker := time.NewTicker(1 * time.Hour)
	quit := make(chan struct{}) // Channel to signal immediate stop
	bondsClient := clients.NewBondsClient()

	defer func() {
		ticker.Stop()
		close(quit)
	}()

	for {
		select {
		case <-ticker.C:
			if b.BondsClientActive {
				result, err := bondsClient.ProcessSavingBondsOffers()
				if err != nil {
					log.Printf("[Bot Service] Failed to get bonds offers: %s", err.Error())
					b.SendMessage(chatId, "There was an error while processing the bonds offers, sorry :(", nil)

				} else if result > 0 {
					timeNow := time.Now()
					message := "<b>12 months savings bonds interest rate has reached the desired value!</b> \n\n" +
						"The current interest rate (" + timeNow.Format("02.01.2006 15:04") + "): <strong>" + fmt.Sprintf("%.2f", result) + "%</strong>\n\n" +
						"<a href='" + b.Config.BondsViewURL + "'>Buy bonds</a>"

					b.SendMessage(chatId, message, nil)
				}
			}
		case <-quit:
			return
		}

		// Break the loop if the bondsClientActive flag is set to false
		if !b.BondsClientActive {
			quit <- struct{}{}
		}
	}
}

func (b *BotService) handleButton(query *tgbotapi.CallbackQuery) {
	var text string

	markup := tgbotapi.NewInlineKeyboardMarkup()
	message := query.Message

	if query.Data == nextButton {
		text = secondMenu
		markup = secondMenuMarkup
	} else if query.Data == backButton {
		text = firstMenu
		markup = firstMenuMarkup
	}

	callbackCfg := tgbotapi.NewCallback(query.ID, "")
	b.Bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	b.Bot.Send(msg)
}

func (b *BotService) SendMessage(chatId int64, text string, entities []tgbotapi.MessageEntity) {
	msg := tgbotapi.NewMessage(chatId, text)
	if len(entities) > 0 {
		msg.Entities = entities
	}
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := b.Bot.Send(msg)

	if err != nil {
		log.Printf("[Bot service] Error sending a message: %s", err.Error())
	}
}

func (b *BotService) SendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := b.Bot.Send(msg)

	return err
}
