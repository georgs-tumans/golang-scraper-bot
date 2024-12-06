package handlers

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"
	"web_scraper_bot/config"
	"web_scraper_bot/utilities"
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
	Code    string
	Ticker  *time.Ticker
	Context context.Context
	Cancel  context.CancelFunc
}

/*
TODO
  - add code uniqueness check
*/
func CreateTracker(code string, runInterval time.Duration) *Tracker {
	ctx, cancel := context.WithCancel(context.Background())

	return &Tracker{
		Code:    code,
		Ticker:  time.NewTicker(runInterval),
		Context: ctx,
		Cancel:  cancel,
	}
}

func (t *Tracker) Stop() {
	t.Ticker.Stop()
	t.Cancel()
}

func (t *Tracker) UpdateInterval(newInterval time.Duration) {
	t.Ticker.Stop()
	t.Ticker = time.NewTicker(newInterval)
}

type CommandHandler struct {
	Config         *config.Configuration
	RunninTrackers []*Tracker
	mu             sync.Mutex
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		Config: config.GetConfig(),
	}
}

func (ch *CommandHandler) HandleCommand(chatId int64, commandString string) error {
	var err error

	// Some commands may have parameters that are separated by a space (/set_interval 5m)
	commandParts := strings.Split(commandString, " ")
	command := strings.ReplaceAll(commandParts[0], "/", "")
	commandParam := ""

	if len(commandParts) > 1 {
		commandParam = commandParts[1]
	}

	/*Commands must always consist of two words separated by an underscore;
	  the first word is the command code, the second word is the actual command.
	  There can be exceptions for commands that target all of the bot functionality instead of specific parts/clients.
	*/
	commandSplit := strings.Split(command, "_")
	commandFunction := command
	commandCode := ""
	if len(commandSplit) > 1 {
		commandCode = commandSplit[0]
		commandFunction = commandSplit[1]
	}

	log.Printf("[CommandHandler] Handling command: %s; parameters: %s", command, commandParam)

	switch commandFunction {

	case "start":
		ch.handleStart(commandCode, chatId)

	case "stop":
		ch.handleStop(commandCode, chatId)

	// case "status":
	// 	bh.handleBondsStatus(chatId)

	case "interval":
		ch.handleSetInterval(commandCode, commandParam, chatId)

	}

	return err
}

// TODO: implement interval setting here
func (ch *CommandHandler) handleStart(code string, chatId int64) {
	if code != "" {
		if tracker := ch.GetActiveTracker(code); tracker == nil {
			go ch.startTracker(code, chatId, 0) // Run in a separate goroutine
			// TODO figure out how to send a message to the user
			// bh.BotFixer.SendMessage(chatId, "Savings bonds client has been started", nil)
			log.Printf("[CommandHandler] Starting tracker: %s", code)
		} else {
			log.Printf("[CommandHandler] Tracker '%s' is already running", code)
			// TODO figure out how to send a message to the user
			//bh.BotFixer.SendMessage(chatId, "Savings bonds client is already running", nil)
		}

	} else {
		/// General /start command that would run all trackers
	}
}

func (ch *CommandHandler) startTracker(code string, chatId int64, runInterval time.Duration) {
	newTracker := CreateTracker(code, runInterval)
	ch.AddRunningTracker(newTracker)

	go func() {
		for {
			select {
			case <-newTracker.Ticker.C:
				// TODO: implement the actual tracker logic here
				// Will call some other struct with logic to make API calls or scrape a website
				// Determine what type of tracker do we start - API or Scraper
				// trackerType := ch.DetermineTrackerType(code)
			case <-newTracker.Context.Done():
				log.Printf("[CommandHandler] Tracker '%s' stopped", code)
				return
			}
		}
	}()
}

func (ch *CommandHandler) handleStop(code string, chatId int64) {
	// Stop all trackers
	if code == "" {
		ch.StopAllTrackers()
	}

	if tracker := ch.GetActiveTracker(code); tracker != nil {
		ch.RemoveRunningTracker(code)
		tracker.Stop()
		// Delete the tracker from the list of running trackers? Is that enough?
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Savings bonds client has been stopped", nil)
	} else {
		log.Printf("[CommandHandler] Tracker '%s' is not running", code)
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Savings bonds client is not running", nil)
	}
}

func (ch *CommandHandler) StopAllTrackers() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.RunninTrackers {
		tracker.Stop()
	}
	ch.RunninTrackers = nil
}

func (ch *CommandHandler) handleSetInterval(code string, newIntervalString string, chatId int64) {
	newInterval, err := utilities.ParseDurationWithDays(newIntervalString)
	if err != nil {
		log.Printf("[CommandHandler] Invalid interval value: %s", err.Error())
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

		return
	}

	if tracker := ch.GetActiveTracker(code); tracker != nil {
		tracker.UpdateInterval(newInterval)
		log.Printf("[CommandHandler] Updated tracker '%s' interval to %s", code, newInterval)
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)
	} else {
		log.Printf("[CommandHandler] Tracker '%s' not found for interval update", code)
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)
	}
}

/******************Utility******************/

func (ch *CommandHandler) DetermineTrackerType(trackerCode string) string {
	for _, tracker := range ch.Config.APITrackers {
		if tracker.Code == trackerCode {
			return API
		}
	}

	for _, tracker := range ch.Config.ScraperTrackers {
		if tracker.Code == trackerCode {
			return Scraper
		}
	}

	return ""
}

func (ch *CommandHandler) GetActiveTracker(trackerCode string) *Tracker {
	// So that only one goroutine can access the trackers at a time
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.RunninTrackers {
		if tracker.Code == trackerCode {
			return tracker
		}
	}

	return nil
}

func (ch *CommandHandler) AddRunningTracker(tracker *Tracker) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.RunninTrackers = append(ch.RunninTrackers, tracker)
}

func (ch *CommandHandler) RemoveRunningTracker(trackerCode string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for i, tracker := range ch.RunninTrackers {
		if tracker.Code == trackerCode {
			ch.RunninTrackers = append(ch.RunninTrackers[:i], ch.RunninTrackers[i+1:]...)
			return
		}
	}
}
