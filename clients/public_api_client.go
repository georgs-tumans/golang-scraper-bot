package clients

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	config "web_scraper_bot/config"
	"web_scraper_bot/helpers"
	"web_scraper_bot/services"

	"github.com/tidwall/gjson"
)

// Client for fetching data from public APIs and extracting the necessary data as defined in the tracker configuration
type PublicAPIClient struct {
	trackerData *config.Tracker
}

func NewPublicAPIClient() *PublicAPIClient {
	return &PublicAPIClient{}
}

func (c *PublicAPIClient) FetchAndExtractData(trackerData *config.Tracker) (*DataResult, error) {
	c.trackerData = trackerData
	dataJson, err := c.getDataFromPublicAPI()
	if err != nil {
		log.Println("[Public API CLient] Error getting data from public API for tracker: "+c.trackerData.Code, err.Error())
		return nil, err
	}

	extractedValue, err := c.extractDataFromPublicAPIResponse(dataJson)
	if err != nil {
		log.Println("[Public API CLient] Error extracting data from public API response for tracker: "+c.trackerData.Code, err.Error())
		return nil, err
	}

	extractedValueFloat, extractedErr := strconv.ParseFloat(extractedValue, 64)
	targetValueFloat, targetErr := strconv.ParseFloat(c.trackerData.NotifyValue, 64)
	if extractedErr != nil || targetErr != nil {
		log.Println("[Public API CLient] Error converting values for tracker: "+c.trackerData.Code, extractedErr.Error(), targetErr.Error())
		return nil, errors.New("error converting values")
	}

	shouldNotify, err := helpers.CompareNumbers(extractedValueFloat, targetValueFloat, c.trackerData.NotifyCriteria)
	if err != nil {
		log.Println("[Public API CLient] Error comparing values for tracker: "+c.trackerData.Code, err.Error())
		return nil, err
	}

	result := &DataResult{
		CurrentValue: extractedValueFloat,
		TargetValue:  targetValueFloat,
		ShouldNotify: shouldNotify,
	}

	c.trackerData = nil

	return result, nil
}

func (c *PublicAPIClient) getDataFromPublicAPI() ([]byte, error) {
	response, err := services.GetRequest(c.trackerData.APIURL)
	if err != nil {
		log.Println("[Public API CLient] Error getting data from public API for tracker: "+c.trackerData.Code, err.Error())

		return nil, err
	}

	return response, nil
}

func (c *PublicAPIClient) extractDataFromPublicAPIResponse(responseJson []byte) (string, error) {
	if responseJson == nil {
		return "", errors.New("nil response data")
	}

	result := gjson.GetBytes(responseJson, c.trackerData.ResponsePath)
	if !result.Exists() {
		log.Println("[Public API CLient] Error extracting data from public API response via the provided JSON path for tracker: "+c.trackerData.Code, "JSON path not found")

		return "", errors.New("json path not found")
	}

	switch result.Value().(type) {
	case string:
		return result.String(), nil
	case float64:
		return fmt.Sprintf("%f", result.Float()), nil
	default:
		log.Println("[Public API CLient] Unrecognized extracted response data type for tracker: "+c.trackerData.Code, "Unsupported data type")
		return "", errors.New("unsupported extracted response data type")
	}
}
