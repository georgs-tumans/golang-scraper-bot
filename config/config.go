package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Configuration struct {
	BondsDataSourceURL string
	BondsViewURL       string
	BondsRateThreshold float64
	BotAPIKey          string
	WebhookURL         string
	Port               string
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
		}

		if rate, rateErr := strconv.ParseFloat(os.Getenv("BONDS_RATE_THRESHOLD"), 64); rateErr != nil {
			log.Fatalf("[GetConfig] Error parsing BONDS_RATE_THRESHOLD")
		} else {
			config.BondsRateThreshold = rate
		}

		// For debugging purposes
		// configJSON, err := json.MarshalIndent(config, "", "  ")
		// if err != nil {
		// 	log.Fatalf("[GetConfig] Error serializing configuration to JSON: %v", err)
		// }
		// log.Printf("[GetConfig] Loaded configuration: %s\n", configJSON)
	}

	return config
}
