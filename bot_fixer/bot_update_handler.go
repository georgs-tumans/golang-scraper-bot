package bot_fixer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"web_scraper_bot/clients"
	"web_scraper_bot/utilities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bondsClient *clients.BondsClient

func (b *BotFixer) webhookHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[Bot fixer] Error reading request body: %v", err)
		http.Error(w, "Could not read request body", http.StatusBadRequest)

		return
	}

	// Parse the body as a Telegram update
	var update tgbotapi.Update
	if err := json.Unmarshal(body, &update); err != nil {
		log.Printf("[Bot fixer] Error parsing update: %v", err)
		http.Error(w, "Could not parse update", http.StatusBadRequest)

		return
	}

	// Handle the update
	b.handleUpdate(update)

	// Respond with a 200 OK status to Telegram
	w.WriteHeader(http.StatusOK)
}

func (b *BotFixer) receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
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

func (b *BotFixer) handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		b.handleMessage(update.Message)

	// Handle button clicks
	case update.CallbackQuery != nil:
		b.handleButton(update.CallbackQuery)
	}
}

func (b *BotFixer) handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	log.Printf("[Bot fixer] %s wrote %s", user.FirstName, text)

	var err error
	if strings.HasPrefix(text, "/") {
		err = b.handleCommand(message.Chat.ID, text)
	}

	if err != nil {
		log.Printf("[Bot fixer] An error occured while handlind message: %s", err.Error())
	}
}

// When we get a command, we react accordingly
func (b *BotFixer) handleCommand(chatId int64, command string) error {
	var err error
	var param string

	// Some commands may have parameters
	commandParts := strings.Split(command, " ")
	mainCommand := commandParts[0]

	if len(commandParts) > 1 {
		param = commandParts[1]
	}

	switch mainCommand {

	case "/menu":
		err = b.SendMenu(chatId)

	case "/delete_webhook":
		if err := b.DeleteWebhook(); err != nil {
			b.SendMessage(chatId, "Failed to delete webhook", nil)
		}

	case "/bonds_start":
		if !b.BondsClientActive {
			b.BondsClientActive = true
			b.SendMessage(chatId, "Savings bonds client has been started", nil)
			go b.activateBondsClient(chatId, 0) // Run in a separate goroutine
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

	case "/bonds_status":
		if b.BondsClientActive {
			var builder strings.Builder
			builder.WriteString("Savings bonds client is currently running\n\n")
			builder.WriteString("Start time: " + bondsClient.ClientStartTimestamp.Format("02.01.2006 15:04") + "\n\n")
			builder.WriteString(bondsClient.FormatOffersMessage())

			b.SendMessage(chatId, builder.String(), nil)
		} else {
			b.SendMessage(chatId, "Savings bonds client is currently not running", nil)
		}

	// TODO: implement update interval validations
	case "/bonds_set_interval":
		if param != "" {
			interval, err := utilities.ParseDurationWithDays(param)
			if err != nil {
				b.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

				return err
			}

			b.BondsClientActive = false
			b.BondsClientActive = true
			go b.activateBondsClient(chatId, interval)
			b.SendMessage(chatId, fmt.Sprintf("Bonds client has been restarted and the run interval updated"), nil)
		} else {
			b.SendMessage(chatId, "Please provide an interval value string in the format <amount><type>.\nExample: 1m, 2h, 1d\nAvailable interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)
		}
	}

	return err
}

func (b *BotFixer) activateBondsClient(chatId int64, runInterval time.Duration) {
	bondsClient = clients.NewBondsClient(runInterval)
	ticker := time.NewTicker(bondsClient.RunInterval)
	quit := make(chan struct{}) // Channel to signal immediate stop

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
					log.Printf("[Bot fixer] Failed to get bonds offers: %s", err.Error())
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

func (b *BotFixer) handleButton(query *tgbotapi.CallbackQuery) {
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
