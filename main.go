package main

import (
	services "bonds_bot/services/bot"
	"io"
	"log"
	"os"
)

func main() {
	logFile := setupLogger()
	defer logFile.Close()

	log.Println("Starting bot service")
	botService := services.NewBotService()
	botService.InitializeBot()
}

func setupLogger() *os.File {
	file, err := os.OpenFile("golang_web_scraper_bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("[setupLogger] Failed to open log file: %v", err)
	}

	// Set the output for the logger to both the file and the console
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return file
}
