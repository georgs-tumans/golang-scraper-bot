package handlers

import (
	"errors"
	"log"
	"strings"
	"sync"
	"web_scraper_bot/config"
	"web_scraper_bot/helpers"
	"web_scraper_bot/utilities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	config          *config.Configuration
	runningTrackers []*Tracker
	commandMap      map[string]CommandFunc
	bot             *tgbotapi.BotAPI
	mu              sync.Mutex
}

type CommandFunc func(code string, chatId int64, commandParam *string) error

func NewCommandHandler(bot *tgbotapi.BotAPI) *CommandHandler {
	ch := &CommandHandler{
		config: config.GetConfig(),
		bot:    bot,
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

	log.Printf("[CommandHandler] Handling command: %s", commandString)

	if handler, exists := ch.commandMap[commandFunction]; exists {
		if err := handler(commandCode, chatId, commandParam); err != nil {
			return err
		}
	} else {
		log.Printf("[CommandHandler] Unknown command: %s", commandFunction)
		helpers.SendMessage(ch.bot, chatId, "Unrecognized command", nil)

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
		newTracker, err := CreateTracker(code, 0, ch.config)
		if err != nil {
			log.Printf("[CommandHandler] Error creating a new tracker: %s", code)
			helpers.SendMessage(ch.bot, chatId, "Error creating a new tracker", nil)
			return err
		}
		ch.AddRunningTracker(newTracker)
		newTracker.Start()

		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' has been started", nil)
		log.Printf("[CommandHandler] Starting tracker: %s", code)
	} else {
		log.Printf("[CommandHandler] Tracker '%s' is already running", code)
		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' is already running", nil)
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
		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' has been stopped", nil)
	} else {
		log.Printf("[CommandHandler] Tracker '%s' is not running", code)
		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' is not running", nil)
	}

	return nil
}

func (ch *CommandHandler) StopAllTrackers() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.runningTrackers {
		tracker.Stop()
	}
	ch.runningTrackers = nil
}

func (ch *CommandHandler) handleSetInterval(code string, chatId int64, commandParam *string) error {
	if commandParam == nil {
		log.Printf("[CommandHandler] No interval value provided")
		helpers.SendMessage(ch.bot, chatId, "No interval value provided", nil)

		return errors.New("no interval value provided")
	}

	newInterval, err := utilities.ParseDurationWithDays(*commandParam)
	if err != nil {
		log.Printf("[CommandHandler] Invalid interval value: %s", err.Error())
		helpers.SendMessage(ch.bot, chatId, "Invalid interval value. Available interval types: 'm'(minute), 'h'(hour), 'd'(day)", nil)

		return err
	}

	if tracker := ch.GetActiveTracker(code); tracker != nil {
		tracker.Stop()
		tracker.UpdateInterval(newInterval)
		tracker.Start()
		log.Printf("[CommandHandler] Updated tracker '%s' interval to %s", code, newInterval)
		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' interval update successfully", nil)

		return nil
	} else {
		log.Printf("[CommandHandler] Tracker '%s' not found for interval update", code)
		helpers.SendMessage(ch.bot, chatId, "Tracker '"+code+"' not found, it's probably not running", nil)

		return errors.New("tracker not found")
	}
}

/******************Utility******************/

func (ch *CommandHandler) GetActiveTracker(trackerCode string) *Tracker {
	// So that only one goroutine can access the trackers at a time
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, tracker := range ch.runningTrackers {
		if tracker.Code == trackerCode {
			return tracker
		}
	}

	return nil
}

func (ch *CommandHandler) AddRunningTracker(tracker *Tracker) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.runningTrackers = append(ch.runningTrackers, tracker)
}

func (ch *CommandHandler) RemoveRunningTracker(trackerCode string) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for i, tracker := range ch.runningTrackers {
		if tracker.Code == trackerCode {
			ch.runningTrackers = append(ch.runningTrackers[:i], ch.runningTrackers[i+1:]...)
			return
		}
	}
}
