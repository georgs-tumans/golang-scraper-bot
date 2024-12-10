package handlers

/*
  TrackerBehavior interface will be implemented by the concrete types of behaviors;
  these behaviors represent different ways of fetching data - either through an API or by scraping a website
  or whatever other way we might come up with in the future. This abstraction is supposed to make it easier
  to add new ways of fetching data without changing the existing code.
  It is NOT meant for implementing the data fetching logic itself - that will be done in the clients.
*/
type TrackerBehavior interface {
	Execute(trackerCode string) error
}

type APITrackerBehavior struct {
	// Client APIClient Client will implement the actual logic of fetching data from the API
}

func (a *APITrackerBehavior) Execute(trackerCode string) error {
	// Call the client for fetching API data and process the result

	// data, err := a.Client.FetchData(a.URL)
	// if err != nil {
	//     return err
	// }

	// log.Printf("[APITracker] Data fetched from API: %v", data)
	return nil
}

type ScraperTrackerBehavior struct {
	// Client ScraperClient Client will implement the actual logic of fetching data by scarping a website
}

func (s *ScraperTrackerBehavior) Execute(trackerCode string) error {
	// Call the client for fetching website data and process the result

	// data, err := s.Client.FetchData(s.URL)
	// if err != nil {
	//     return err
	// }

	// log.Printf("[ScraperTracker] Data scraped from website: %v", data)
	return nil
}
