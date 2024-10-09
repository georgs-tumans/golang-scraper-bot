package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Configuration struct {
	SMPTHost           string
	SMTPPort           int
	SMTPPassword       string
	SMTPUsername       string
	EmailRecipients    []string
	BondsDataSourceURL string
	BondsViewURL       string
	BondsRateThreshold float64
	BotAPIKey          string
}

var config *Configuration

func GetConfig() *Configuration {
	if config == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("[GetConfig] Error loading .env file")
		}

		config = &Configuration{
			SMPTHost:           os.Getenv("SMTP_HOST"),
			SMTPPassword:       os.Getenv("SMTP_PASSWORD"),
			SMTPUsername:       os.Getenv("SMTP_USERNAME"),
			BondsDataSourceURL: os.Getenv("BONDS_DATA_SOURCE_URL"),
			BondsViewURL:       os.Getenv("BONDS_VIEW_URL"),
			BotAPIKey:          os.Getenv("BOT_API_KEY"),
		}

		if port, err := strconv.Atoi(os.Getenv("SMTP_PORT")); err != nil {
			log.Fatalf("[GetConfig] Error parsing SMTP_PORT")
		} else {
			config.SMTPPort = port
		}

		if rate, rateErr := strconv.ParseFloat(os.Getenv("BONDS_RATE_THRESHOLD"), 64); rateErr != nil {
			log.Fatalf("[GetConfig] Error parsing BONDS_RATE_THRESHOLD")
		} else {
			config.BondsRateThreshold = rate
		}

		recipients := os.Getenv("EMAIL_RECIPIENTS")
		config.EmailRecipients = strings.Split(recipients, ",")
	}

	return config
}
