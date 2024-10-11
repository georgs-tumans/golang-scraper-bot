# Web Scraper Telegram Bot

## About

A Telegram bot with a set of tools originally intended for acquiring whatever data neccessary from the web and notifying the user upon reaching certain criteria.

## Available tools/functionality

 - Fetching data about the current interest rates of the government issued savings bonds of the Republic of Latvia notifying the user in case the 12 months bonds interest rate is equal or higher than the desired configured value.

## Preconditions

- A Telegram bot API key which means you must register a bot. Learn how to do it [here](https://core.telegram.org/bots#how-do-i-create-a-bot).
- Docker installed (if you want to run this in a container)
- GO installed (if you want to run it as a regular console app)

## Use

1. Register a bot with Telegram
2. Build and run this app
3. Use commands to interact with your new Telegram bot :)


## Initial setup

Before running, you **must** create an `.env` file; use this [example](/.env.example) to fill out the values.

Then use [this](/docker_build_and_run.ps1) included powershell script to build and run the bot as a Docker container.

Afterwards you can use [this other script](/docker_build_and_run.ps1) to run the container without rebuilding the image.

Also, you can press `F5` if using VS Code to run via a launch profile or just use the CMD command `go run main.go` in the root of the project.