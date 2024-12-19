package clients

import (
	"log"
	"strconv"
	config "web_scraper_bot/config"
	"web_scraper_bot/helpers"
)

type DataResult struct {
	CurrentValue float64
	ShouldNotify bool
}

// Common client interface that will be implemented by the concrete types of clients
type Client interface {
	FetchAndExtractData(trackerCode string) (*DataResult, error)
}

func ShouldNotify(trackerData *config.Tracker, extractedValue float64) (bool, error) {
	for _, criteria := range trackerData.NotifyCriteria {
		notifyValue, err := strconv.ParseFloat(criteria.Value, 64)
		if err != nil {
			log.Println("[Client] Error converting notification criteria value for tracker: "+trackerData.Code, err.Error())
			return false, err
		}

		result, err := helpers.CompareNumbers(extractedValue, notifyValue, criteria.Operator)
		if err != nil {
			log.Println("[Client] Error comparing extracted and notification target values for tracker: "+trackerData.Code, err.Error())
			return false, err
		}

		return result, nil
	}

	return false, nil
}
