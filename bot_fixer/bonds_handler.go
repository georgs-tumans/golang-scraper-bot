package bot_fixer

import (
	"fmt"
	"log"
	"strings"
	"time"
	clients "web_scraper_bot/clients"
	"web_scraper_bot/utilities"
)

type BondsHandler struct {
	BotFixer          *BotFixer
	BondsClientActive bool
	BondsClient       *clients.BondsClient
}

func NewBondsHandler(botFixer *BotFixer) *BondsHandler {
	return &BondsHandler{
		BondsClientActive: false,
		BotFixer:          botFixer,
	}
}

func (bh *BondsHandler) HandleBondsCommand(chatId int64, command string) error {
	var err error

	// Some commands may have parameters
	commandParts := strings.Split(command, " ")
	mainCommand := commandParts[0]
	param := ""

	if len(commandParts) > 1 {
		param = commandParts[1]
	}

	log.Printf("[BondsHandler] Handling command: %s; parameters: %s", mainCommand, param)

	switch mainCommand {

	case "/bonds_start":
		bh.handleBondsStart(chatId)

	case "/bonds_stop":
		bh.handleBondsStop(chatId)

	case "/bonds_status":
		bh.handleBondsStatus(chatId)

	case "/bonds_set_interval":
		err = bh.handleBondsSetInterval(chatId, param)
	}

	return err

}

func (bh *BondsHandler) handleBondsStart(chatId int64) {
	if !bh.BondsClientActive {
		bh.BondsClientActive = true
		bh.BotFixer.SendMessage(chatId, "Savings bonds client has been started", nil)
		log.Printf("[BondsHandler] Starting the bonds client")
		go bh.activateBondsClient(chatId, 0) // Run in a separate goroutine
	} else {
		log.Printf("[BondsHandler] Bonds client is already running")
		bh.BotFixer.SendMessage(chatId, "Savings bonds client is already running", nil)
	}
}

func (bh *BondsHandler) activateBondsClient(chatId int64, runInterval time.Duration) {
	if bh.BondsClient != nil && bh.BondsClient.Ticker != nil {
		log.Println("[BondsHandler] Stopping existing ticker")
		bh.BondsClient.StopTicker()
	}

	bh.BondsClient = clients.NewBondsClient(runInterval)
	quit := make(chan struct{}) // Channel to signal immediate stop

	go func() {
		for {
			select {
			case <-bh.BondsClient.Ticker.C:
				if bh.BondsClientActive {
					result, err := bh.BondsClient.ProcessSavingBondsOffers()
					if err != nil {
						log.Printf("[BondsHandler] Error while processing bonds offers: %s", err.Error())
						bh.BotFixer.SendMessage(chatId, "Error processing bonds offers", nil)
					} else if result > 0 {
						log.Println("[BondsHandler] Notifying the user about the desired interest rate")
						timeNow := time.Now()

						var builder strings.Builder
						builder.WriteString("<b>12 months savings bonds interest rate has reached the desired value!</b> \n\n")
						builder.WriteString(fmt.Sprintf("The current interest rate (%s): <strong> %.2f%% </strong>\n\n", timeNow.Format("02.01.2006 15:04"), result))
						builder.WriteString(fmt.Sprintf("<a href='%s'> View in browser </a>", bh.BotFixer.Config.BondsViewURL))

						bh.BotFixer.SendMessage(chatId, builder.String(), nil)
					}
				}
			case <-quit:
				log.Println("[BondsHandler] Stopping bonds client")
				bh.BondsClient.StopTicker()

				return
			}
		}
	}()
}

func (bh *BondsHandler) handleBondsStop(chatId int64) {
	if bh.BondsClientActive {
		bh.BondsClientActive = false
		bh.BotFixer.SendMessage(chatId, "Savings bonds client has been stopped", nil)
	} else {
		bh.BotFixer.SendMessage(chatId, "Savings bonds client is not running", nil)
	}
}

func (bh *BondsHandler) handleBondsStatus(chatId int64) {
	if bh.BondsClientActive {
		var builder strings.Builder
		builder.WriteString("Savings bonds client is currently running\n\n")
		builder.WriteString("Start time: " + bh.BondsClient.ClientStartTimestamp.Format("02.01.2006 15:04") + "\n\n")
		builder.WriteString("Current run interval: " + bh.BondsClient.RunInterval.String() + "\n\n")
		builder.WriteString(bh.BondsClient.FormatOffersMessage())
		bh.BotFixer.SendMessage(chatId, builder.String(), nil)
	} else {
		bh.BotFixer.SendMessage(chatId, "Savings bonds client is currently not running", nil)
	}
}

func (bh *BondsHandler) handleBondsSetInterval(chatId int64, param string) error {
	if param != "" {
		interval, err := utilities.ParseDurationWithDays(param)
		if err != nil {
			log.Printf("[BondsHandler] Invalid interval value: %s", err.Error())
			bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

			return err
		}

		bh.BondsClientActive = false
		time.Sleep(2 * time.Second) // Wait for the client to stop
		bh.BondsClientActive = true
		bh.BotFixer.SendMessage(chatId, "Bonds client has been restarted and the run interval updated", nil)
		go bh.activateBondsClient(chatId, interval)
	} else {
		var builder strings.Builder
		builder.WriteString("Invalid command: missing the interval value\n\n")
		builder.WriteString("Correct command use example:\n <code>/bonds_set_interval 5m</code>\n\n")
		builder.WriteString("Available interval types: 'm'(minute), 'h'(hour), 'd'(day)")
		bh.BotFixer.SendMessage(chatId, builder.String(), nil)
		log.Printf("[BondsHandler] Interval value required")

		return fmt.Errorf("interval_value_required")
	}

	return nil
}
