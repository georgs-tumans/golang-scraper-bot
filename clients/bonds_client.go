package clients

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"web_scraper_bot/config"
	"web_scraper_bot/services"
	"web_scraper_bot/utilities"
)

type BondsClient struct {
	BondsDataSourceURL   string
	BondsViewURL         string
	BondsRateThreshold   float64
	ClientStartTimestamp time.Time
	LastRunTimestamp     time.Time
	LastBondsOffers      *OffersResponse
	RunInterval          time.Duration
	Ticker               *time.Ticker
}

type Offer struct {
	InterestRate float64 `json:"interestRate"`
	Period       int     `json:"period"`
}

type OffersResponse []*Offer

func NewBondsClient(runInterval time.Duration) *BondsClient {
	config := config.GetConfig()

	newClient := &BondsClient{
		BondsDataSourceURL:   config.BondsDataSourceURL,
		BondsViewURL:         config.BondsViewURL,
		BondsRateThreshold:   config.BondsRateThreshold,
		ClientStartTimestamp: time.Now(),
		RunInterval:          runInterval,
	}

	if runInterval == 0 {
		runInterval, err := utilities.ParseDurationWithDays(config.BondsRunInterval)
		if err != nil {
			log.Fatalf("[NewBondsClient] Failed to parse bonds client run interval: %v", err)
		}
		newClient.RunInterval = runInterval
	}

	newClient.Ticker = time.NewTicker(newClient.RunInterval)

	return newClient
}

func (c *BondsClient) getBondsOffers() (*OffersResponse, error) {
	offersResponse := &OffersResponse{}
	if err := services.GetRequest(c.BondsDataSourceURL, offersResponse); err != nil {
		log.Println("[getBondsOffers] Failed to get data")

		return nil, err
	}

	c.LastBondsOffers = offersResponse
	c.LastRunTimestamp = time.Now()
	responseJsonString, err := json.Marshal(offersResponse)
	if err != nil {
		log.Println("[getBondsOffers] Failed to marshal response", err)
	} else {
		log.Println("[getBondsOffers] Current offers", string(responseJsonString))
	}

	return offersResponse, nil
}

func (c *BondsClient) ProcessSavingBondsOffers() (float64, error) {
	log.Println("[ProcessSavingBondsOffers] Processing saving bonds offers")
	bondOffers, err := c.getBondsOffers()
	if err != nil {
		return 0, err
	}

	for _, offer := range *bondOffers {
		if offer.Period == 12 && offer.InterestRate >= c.BondsRateThreshold {
			log.Println("[ProcessSavingBondsOffers] 12 months interest rate match (" + fmt.Sprintf("%.2f", offer.InterestRate) + ")")
			return offer.InterestRate, nil
		}
	}

	return 0, nil
}

func (c *BondsClient) FormatOffersMessage() string {
	if c.LastBondsOffers == nil {
		return ""
	}

	var builder strings.Builder
	builder.WriteString("Latest updates (" + c.LastRunTimestamp.Format("02.01.2006 15:04") + "):\n")
	for _, offer := range *c.LastBondsOffers {
		builder.WriteString(fmt.Sprintf("Period: %d months, Interest rate: %.2f%%\n", offer.Period, offer.InterestRate))
	}

	return builder.String()
}

func (c *BondsClient) StopTicker() {
	if c.Ticker != nil {
		c.Ticker.Stop()
		c.Ticker = nil
	}
}
