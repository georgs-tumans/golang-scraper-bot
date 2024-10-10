package clients

import (
	"bonds_bot/config"
	"bonds_bot/services"
	"encoding/json"
	"fmt"
	"log"
)

type BondsClient struct {
	BondsDataSourceURL string
	BondsViewURL       string
	BondsRateThreshold float64
}

type Offer struct {
	InterestRate float64 `json:"interestRate"`
	Period       int     `json:"period"`
}

type OffersResponse []*Offer

func NewBondsClient() *BondsClient {
	config := config.GetConfig()

	return &BondsClient{
		BondsDataSourceURL: config.BondsDataSourceURL,
		BondsViewURL:       config.BondsViewURL,
		BondsRateThreshold: config.BondsRateThreshold,
	}
}

func (c *BondsClient) getBondsOffers() (*OffersResponse, error) {
	offersResponse := &OffersResponse{}
	if err := services.GetRequest(c.BondsDataSourceURL, offersResponse); err != nil {
		log.Println("[getBondsOffers] Failed to get data")

		return nil, err
	}

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
