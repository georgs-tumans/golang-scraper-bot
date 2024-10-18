# Web Scraper Telegram Bot

## About

A Telegram bot with a set of tools originally intended for acquiring whatever data neccessary from the web and notifying the user upon reaching certain criteria.

## Available tools/functionality

 - Fetching data about the current interest rates of the government issued savings bonds of the Republic of Latvia notifying the user in case the 12 months bonds interest rate is equal or higher than the desired configured value.

## Preconditions

- A Telegram bot API key which means you must register a bot. Learn how to do it [here](https://core.telegram.org/bots#how-do-i-create-a-bot).
- Docker installed (if you want to run this in a container)
- GO installed (if you want to run it as a regular console app)
- ngrok running for local development

## Use

1. Register a bot with Telegram
2. Build and run this app
3. Use commands to interact with your new Telegram bot :)


## Initial setup

1. Run ngrok locally - you will need it for exposing localhost to the internet so that Telegram can reach the bot when running locally (during development). 

There is a [powershell script](/docker_run_ngrok.ps1) for hassle free setup of ngrok via Docker but in order to use it:

* Create an ngrok configuration file `ngrok.yml` based on this [template](./ngrok.yml.example)
* Edit the [script](/docker_run_ngrok.ps1) and set the location of the newly created `ngrok.yml`
* Run the script

You will need to know the ngrok generated URL that tunnels your locally run app to the internet - open `http://localhost:4040/status` in a browser to view the ngrok panel

2. Create an `.env` file; use this [example](/.env.example) to fill out the values.

3. Use [this](/docker_build_and_run.ps1) included powershell script to build (or rebuild) and run the bot as a Docker container.

Afterwards you can use [this other script](/docker_run.ps1) to run the container without rebuilding the image.

Also, you can press `F5` if using VS Code to run via a launch profile or just use the CMD command `go run main.go` in the root of the project.


## Development

In order to develop and run the bot locally via you IDE you must set the environmental variable `ENVIROMENT` to `local`. When running the bot, this will result in automatic deletion of any webhooks registered for the given Telegram bot API key and switching to the long polling approach which works much better for local development.

However, if you need to run the bot in a container or host it somewhere, it is recommended to set the `ENVIROMENT` to `cloud`/`docker` which will the register a new webhook upon instantiation and use that for getting updates.

**NB**

You cannot run the bot in the long polling mode while there are actively registered webhooks for the same bot API key!

## Deployment

The bot is currently hosted on Google Apps; pushing to `master` triggers a build.
