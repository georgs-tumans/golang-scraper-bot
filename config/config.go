package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Tracker struct {
	Code         string  `json:"code"`
	URL          string  `json:"url"`
	Interval     string  `json:"interval"`
	NotifyValue  float64 `json:"notifyValue"`
	ResponsePath string  `json:"responsePath"`
}

type Configuration struct {
	BondsDataSourceURL string
	BondsViewURL       string
	BondsRateThreshold float64
	BotAPIKey          string
	WebhookURL         string
	Port               string
	Environment        string
	BondsRunInterval   string
	APITrackers        []*Tracker
	ScraperTrackers    []*Tracker
}

var config *Configuration

func GetConfig() *Configuration {
	if config == nil {
		err := godotenv.Load()
		if err != nil {
			log.Println("[GetConfig] Error loading .env file")
		}

		config = &Configuration{

			BondsDataSourceURL: os.Getenv("BONDS_DATA_SOURCE_URL"),
			BondsViewURL:       os.Getenv("BONDS_VIEW_URL"),
			BotAPIKey:          os.Getenv("BOT_API_KEY"),
			WebhookURL:         os.Getenv("WEBHOOK_URL"),
			Port:               os.Getenv("PORT"),
			Environment:        os.Getenv("ENVIRONMENT"),
			BondsRunInterval:   os.Getenv("BONDS_RUN_INTERVAL"),
		}

		if rate, rateErr := strconv.ParseFloat(os.Getenv("BONDS_RATE_THRESHOLD"), 64); rateErr != nil {
			log.Fatalf("[GetConfig] Error parsing BONDS_RATE_THRESHOLD")
		} else {
			config.BondsRateThreshold = rate
		}

		if config.BondsRunInterval == "" {
			config.BondsRunInterval = "1h"
		}

		// TODO: Load trackers from the environment

		// For debugging purposes
		// configJSON, err := json.MarshalIndent(config, "", "  ")
		// if err != nil {
		// 	log.Fatalf("[GetConfig] Error serializing configuration to JSON: %v", err)
		// }
		// log.Printf("[GetConfig] Loaded configuration: %s\n", configJSON)
	}

	return config
}
