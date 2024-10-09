package main

import (
	clients "bonds_bot/clients"
	config "bonds_bot/config"
	"bufio"
	"context"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	// Menu texts
	firstMenu  = "<b>Menu 1</b>\n\nA beautiful menu with a shiny inline button."
	secondMenu = "<b>Menu 2</b>\n\nA better menu with even more shiny inline buttons."

	// Button texts
	nextButton     = "Next"
	backButton     = "Back"
	tutorialButton = "Tutorial"

	bot *tgbotapi.BotAPI

	// Client activity statuses
	bondsClientActive = false

	// Client instances
	bondsClient = clients.NewBondsClient()

	// Keyboard layout for the first menu. One button, one row
	firstMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(nextButton, nextButton),
		),
	)

	// Keyboard layout for the second menu. Two buttons, one per row
	secondMenuMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(tutorialButton, "https://core.telegram.org/bots/api"),
		),
	)
)

func main() {
	var err error
	config := config.GetConfig()

	bot, err = tgbotapi.NewBotAPI(config.BotAPIKey)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)

	// Handle button clicks
	case update.CallbackQuery != nil:
		handleButton(update.CallbackQuery)
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("%s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func handleCommand(chatId int64, command string) error {
	var err error

	switch command {

	case "/menu":
		err = sendMenu(chatId)

	case "/bonds_start":
		if !bondsClientActive {
			bondsClientActive = true
			sendMessage(chatId, "Running the savings bonds client")
			go activateBondsClient(chatId) // Run in a separate goroutine
		} else {
			sendMessage(chatId, "Bonds client is already running")
		}

	case "/bonds_stop":
		if bondsClientActive {
			bondsClientActive = false
			sendMessage(chatId, "Stopped the savings bonds client")
		} else {
			sendMessage(chatId, "Bonds client is not running")
		}
	}

	return err
}

func handleButton(query *tgbotapi.CallbackQuery) {
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
	bot.Send(callbackCfg)

	// Replace menu text and keyboard
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	bot.Send(msg)
}

func sendMenu(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, firstMenu)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = firstMenuMarkup
	_, err := bot.Send(msg)

	return err
}

func sendMessage(chatId int64, text string) {
	msg := tgbotapi.NewMessage(chatId, text)
	_, err := bot.Send(msg)

	if err != nil {
		log.Printf("Error sending a message: %s", err.Error())
	}
}

func activateBondsClient(chatId int64) {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{}) // Channel to signal immediate stop

	defer func() {
		ticker.Stop()
		close(quit)
	}()

	for {
		select {
		case <-ticker.C:
			if bondsClientActive {
				sendMessage(chatId, "Pushing notification..")
			}
		case <-quit:
			return
		}

		// Break the loop if the bondsClientActive flag is set to false
		if !bondsClientActive {
			quit <- struct{}{}
		}
	}
}
