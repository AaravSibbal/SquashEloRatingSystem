package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	bot "github.com/AaravSibbal/SqashEloRatingSystem/Bot"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

	func main() {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		
		fmt.Printf("%v", err)
		return	
	}	
	
	db, err := getDbConnection(envFile)
	if err != nil {
		fmt.Printf("%v",err)
		return
	}
	
	
	botToken := envFile["BOT_TOKEN"]
	ctx := context.Background()
	botStruct := &bot.Bot{
		BotToken: botToken,
		Db: db,
		Ctx: &ctx,
	}


	fmt.Println(envFile["BOT_TOKEN"])
	
	botStruct.Run()
}

func getDbConnection(envFile map[string]string) (*sql.DB, error){

	host := envFile["POSTGRES_HOST"]
	postPort := envFile["POSTGRES_PORT"]
	user := envFile["POSTGRES_USER"]
	password := envFile["POSTGRES_PASSWORD"]
	dbName := envFile["POSTGRES_NAME"]

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, postPort, user, password, dbName)

	fmt.Println(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatalln("there was an error getting the db connection", err)
		return nil, err
	}


	if err = db.Ping(); err != nil {
		log.Fatalln("we couldn't ping the db for some reason", err)
	}

	fmt.Println("db was connected successfuly")

	return db, nil
}