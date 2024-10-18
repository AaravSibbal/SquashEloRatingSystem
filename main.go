package main

import (
	"fmt"

	bot "github.com/AaravSibbal/SqashEloRatingSystem/Bot"
	"github.com/joho/godotenv"
)

func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		fmt.Errorf("there was an error reading the .env file,\n\n %w", err)
		return
	}	
	bot.BotToken = envFile["BotToken"]
	fmt.Println(envFile["BotToken"])
	
	bot.Run()
}