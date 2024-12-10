package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	config "web_scraper_bot/config"
	"web_scraper_bot/helpers"
	"web_scraper_bot/services"

	"github.com/tidwall/gjson"
)

func ProcessPublicAPICall(trackerCode string) {
	configuration := config.GetConfig()
	trackerData := configuration.GetAPITrackerData(trackerCode)

	dataJson, err := getDataFromPublicAPI(trackerData)
	if err != nil {
		log.Println("[Public API CLient] Error getting data from public API for tracker: "+trackerCode, err.Error())
		// TODO: How do we handle this, do we tell the user about each failed call? Do some retries or save failure stats?

		return
	}

	extractedValue, err := extractDataFromPublicAPIResponse(trackerData, dataJson)
	if err != nil {
		log.Println("[Public API CLient] Error extracting data from public API response for tracker: "+trackerCode, err.Error())
		// TODO: Same questtion about error handling
		return
	}

	// TODO: make a conversion helper
	shouldNotify, err := helpers.NumberComparison(extractedValue, trackerData.NotifyValue, trackerData.NotifyCriteria)

}

func getDataFromPublicAPI(trackerData *config.Tracker) ([]byte, error) {

	var response interface{}

	if err := services.GetRequest(trackerData.APIURL, response); err != nil {
		log.Println("[Public API CLient] Error getting data from public API for tracker: "+trackerData.Code, err.Error())

		return nil, err
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Println("[Public API CLient] Error marshalling public API response for tracker: "+trackerData.Code, err.Error())

		return nil, err
	}
	log.Println("[Public API CLient] Marhsalled public API response", string(responseJson))

	return responseJson, nil
}

func extractDataFromPublicAPIResponse(trackerData *config.Tracker, responseJson []byte) (string, error) {
	if responseJson == nil {
		return "", errors.New("nil response data")
	}

	result := gjson.GetBytes(responseJson, trackerData.ResponsePath)
	if !result.Exists() {
		log.Println("[Public API CLient] Error extracting data from public API response via the provided JSON path for tracker: "+trackerData.Code, "JSON path not found")

		return "", errors.New("json path not found")
	}

	switch result.Value().(type) {
	case string:
		return result.String(), nil
	case float64:
		return fmt.Sprintf("%f", result.Float()), nil
	default:
		log.Println("[Public API CLient] Unrecognized extracted response data type for tracker: "+trackerData.Code, "Unsupported data type")
		return "", errors.New("unsupported extracted response data type")
	}
}
