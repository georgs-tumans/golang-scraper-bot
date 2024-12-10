package handlers

import (
	"context"
	"fmt"
	"log"
	"time"
	"web_scraper_bot/config"
)

const (
	API     = "api"
	Scraper = "scraper"
)

// Tracker represents a single URL that the bot will track - either through an API or by scraping a website
/* TODO
   - think about reusing one struct for Config and here
*/
type Tracker struct {
	Code     string
	Ticker   *time.Ticker
	Context  context.Context
	Cancel   context.CancelFunc
	Behavior TrackerBehavior
	running  bool
}

/*
TODO
  - add code uniqueness check
*/
func CreateTracker(code string, runInterval time.Duration, config *config.Configuration) (*Tracker, error) {
	var behavior TrackerBehavior
	trackerType := DetermineTrackerType(code, config)

	switch trackerType {
	case API:
		behavior = &APITrackerBehavior{}
	case Scraper:
		behavior = &ScraperTrackerBehavior{}
	default:
		return nil, fmt.Errorf("unsupported client type for code: %s", code)
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Tracker{
		Code:     code,
		Ticker:   time.NewTicker(runInterval),
		Context:  ctx,
		Cancel:   cancel,
		Behavior: behavior,
	}, nil
}

func (t *Tracker) Start() {
	if t.running {
		return
	}
	t.running = true

	go func() {
		defer func() { t.running = false }()
		for {
			select {
			case <-t.Ticker.C:
				if err := t.Behavior.Execute(t.Code); err != nil {
					log.Printf("[Tracker] Error executing tracker '%s': %s", t.Code, err)
				}
			case <-t.Context.Done():
				log.Printf("[Tracker] Stopping tracker '%s'", t.Code)
				return
			}
		}
	}()
}

func (t *Tracker) Stop() {
	if !t.running {
		return
	}

	t.Ticker.Stop()
	t.Cancel()
}

func (t *Tracker) UpdateInterval(newInterval time.Duration) {
	t.Ticker.Stop()
	t.Ticker = time.NewTicker(newInterval)
}

func DetermineTrackerType(trackerCode string, config *config.Configuration) string {
	for _, tracker := range config.APITrackers {
		if tracker.Code == trackerCode {
			return API
		}
	}

	for _, tracker := range config.ScraperTrackers {
		if tracker.Code == trackerCode {
			return Scraper
		}
	}

	return ""
}
