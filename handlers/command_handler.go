package handlers

import (
	"errors"
	"log"
	"strings"
	"sync"
	"web_scraper_bot/config"
	"web_scraper_bot/utilities"
)

type CommandHandler struct {
	Config          *config.Configuration
	RunningTrackers []*Tracker
	commandMap      map[string]CommandFunc
	mu              sync.Mutex
}

type CommandFunc func(code string, chatId int64, commandParam *string) error

func NewCommandHandler() *CommandHandler {
	ch := &CommandHandler{
		Config: config.GetConfig(),
	}

	// TODO: add status command
	ch.commandMap = map[string]CommandFunc{
		"start":    ch.handleStart,
		"stop":     ch.handleStop,
		"interval": ch.handleSetInterval,
	}

	return ch
}

func (ch *CommandHandler) HandleCommand(chatId int64, commandString string) error {
	// Some commands may have parameters that are separated by a space (/set_interval 5m)
	commandParts := strings.Split(commandString, " ")
	command := strings.ReplaceAll(commandParts[0], "/", "")
	var commandParam *string

	if len(commandParts) > 1 {
		commandParam = &commandParts[1]
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

	log.Printf("[CommandHandler] Handling command: %s; parameters: %s", command, *commandParam)

	if handler, exists := ch.commandMap[commandFunction]; exists {
		if err := handler(commandCode, chatId, commandParam); err != nil {
			return err
		}
	} else {
		log.Printf("[CommandHandler] Unknown command: %s", commandFunction)
		// TODO: Send a message to the user indicating an unknown command.
		// bh.BotFixer.SendMessage(chatId, "Error creating a new tracker", nil)

		return errors.New("unknown command")
	}

	return nil
}

// TODO: implement interval setting here
func (ch *CommandHandler) handleStart(code string, chatId int64, commandParam *string) error {
	if code == "" {
		/// General /start command that would run all trackers
		return nil
	}

	if tracker := ch.GetActiveTracker(code); tracker == nil {
		newTracker, err := CreateTracker(code, 0, ch.Config)
		if err != nil {
			log.Printf("[CommandHandler] Error creating a new tracker: %s", code)
			// TODO figure out how to send a message to the user
			// bh.BotFixer.SendMessage(chatId, "Error creating a new tracker", nil)
			return err
		}
		ch.AddRunningTracker(newTracker)
		newTracker.Start()

		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Savings bonds client has been started", nil)
		log.Printf("[CommandHandler] Starting tracker: %s", code)
	} else {
		log.Printf("[CommandHandler] Tracker '%s' is already running", code)
		// TODO figure out how to send a message to the user
		//bh.BotFixer.SendMessage(chatId, "Savings bonds client is already running", nil)
	}

	return nil
}

func (ch *CommandHandler) handleStop(code string, chatId int64, commandParam *string) error {
	// Stop all trackers
	if code == "" {
		ch.StopAllTrackers()

		return nil
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

	return nil
}

func (ch *CommandHandler) StopAllTrackers() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.RunningTrackers {
		tracker.Stop()
	}
	ch.RunningTrackers = nil
}

func (ch *CommandHandler) handleSetInterval(code string, chatId int64, commandParam *string) error {
	if commandParam == nil {
		log.Printf("[CommandHandler] No interval value provided")
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "No interval value provided", nil)

		return errors.New("no interval value provided")
	}

	newInterval, err := utilities.ParseDurationWithDays(*commandParam)
	if err != nil {
		log.Printf("[CommandHandler] Invalid interval value: %s", err.Error())
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

		return err
	}

	if tracker := ch.GetActiveTracker(code); tracker != nil {
		tracker.Stop()
		tracker.UpdateInterval(newInterval)
		tracker.Start()
		log.Printf("[CommandHandler] Updated tracker '%s' interval to %s", code, newInterval)
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

		return nil
	} else {
		log.Printf("[CommandHandler] Tracker '%s' not found for interval update", code)
		// TODO figure out how to send a message to the user
		// bh.BotFixer.SendMessage(chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

		return errors.New("tracker not found")
	}
}

/******************Utility******************/

func (ch *CommandHandler) GetActiveTracker(trackerCode string) *Tracker {
	// So that only one goroutine can access the trackers at a time
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.RunningTrackers {
		if tracker.Code == trackerCode {
			return tracker
		}
	}

	return nil
}

func (ch *CommandHandler) AddRunningTracker(tracker *Tracker) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.RunningTrackers = append(ch.RunningTrackers, tracker)
}

func (ch *CommandHandler) RemoveRunningTracker(trackerCode string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for i, tracker := range ch.RunningTrackers {
		if tracker.Code == trackerCode {
			ch.RunningTrackers = append(ch.RunningTrackers[:i], ch.RunningTrackers[i+1:]...)
			return
		}
	}
}
