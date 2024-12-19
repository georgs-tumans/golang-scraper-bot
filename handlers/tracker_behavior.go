package handlers

import (
	"fmt"
	"web_scraper_bot/clients"
	"web_scraper_bot/config"
	"web_scraper_bot/helpers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
TrackerBehavior interface will be implemented by the concrete types of behaviors;
these behaviors represent different ways of fetching data - either through an API or by scraping a website
or whatever other way we might come up with in the future. This abstraction is supposed to make it easier
to add new ways of fetching data without changing the existing code.
It is NOT meant for implementing the data fetching logic itself - that will be done in the clients.
*/
type TrackerBehavior interface {
	Execute(trackerData *config.Tracker, chatID int64) (string, error)
}

type APITrackerBehavior struct {
	bot    *tgbotapi.BotAPI
	client *clients.PublicAPIClient
}

func NewAPITrackerBehavior(bot *tgbotapi.BotAPI) *APITrackerBehavior {
	return &APITrackerBehavior{
		bot:    bot,
		client: clients.NewPublicAPIClient(),
	}
}

func (tb *APITrackerBehavior) Execute(trackerData *config.Tracker, chatID int64) (string, error) {
	result, err := tb.client.FetchAndExtractData(trackerData)
	if err != nil {
		// Notify the user? Add to some failure statistics?
		return "", err
	}

	// TODO add proper message
	// Rethink sending messages since multiple notification criteria can be set now
	if result.ShouldNotify {
		helpers.SendMessageHTML(tb.bot, chatID, "Notify user about the data", nil)
	}

	return fmt.Sprintf("%.2f", result.CurrentValue), nil
}

type ScraperTrackerBehavior struct {
	bot *tgbotapi.BotAPI
}

func NewScraperTrackerBehavior(bot *tgbotapi.BotAPI) *ScraperTrackerBehavior {
	return &ScraperTrackerBehavior{
		bot: bot,
	}
}

func (s *ScraperTrackerBehavior) Execute(trackerData *config.Tracker, chatID int64) (string, error) {
	// Call the client for fetching website data and process the result

	// data, err := s.Client.FetchData(s.URL)
	// if err != nil {
	//     return err
	// }

	// log.Printf("[ScraperTracker] Data scraped from website: %v", data)
	return "", nil
}
