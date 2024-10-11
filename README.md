# Web Scraper Telegram Bot

## About

A Telegram bot with a set of tools originally intended for acquiring whatever data neccessary from the web and notifying the user upon reaching certain criteria.

## Available tools/functionality

 - Fetching data about the current interest rates of the government issued savings bonds of the Republic of Latvia notifying the user in case the 12 months bonds interest rate is equal or higher than the desired configured value.

## Initial setup

Before running, you **must** create an `.env` file; use this [example](/.env.example) to fill out the values.

Then use [this](/docker_build_and_run.ps1) included powershell script to build and run the bot as a Docker container.

Afterwards you can use [this other script](/docker_build_and_run.ps1) to run the container without rebuilding the image.

Also, you can press `F5` if using VS Code to run via a launch profile or just use the CMD command `go run main.go` in the root of the project.
 
## Use

Intended to be run regularly with the help of automatization tools or something like a Windows Scheduler.

You can use `go build` in the project root to create an .exe file (on Windows).