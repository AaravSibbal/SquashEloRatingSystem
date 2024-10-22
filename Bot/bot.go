package bot

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	BotToken string
	Db *sql.DB
	Ctx *context.Context
}


func (b *Bot) Run() {
	discord, err := discordgo.New("Bot " + b.BotToken)
	if err != nil {
		log.Fatal("There was an error creating a discord bot")
	}

	discord.AddHandler(b.newMessage)

	discord.Open()
	defer discord.Close()

	fmt.Println("The Bot is running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (b *Bot) newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	fmt.Println("we are getting here")
	// if message.Author.ID == discord.State.User.ID {
	// 	return
	// }


	fmt.Printf("Message Content: %s\n", message.Content)
	switch {
	case strings.Contains(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, getHelpMessage())
	case strings.Contains(message.Content, "!ping"):
		discord.ChannelMessageSend(message.ChannelID, ping())
	// default:
	// 	discord.ChannelMessageSend(message.ChannelID, "something went wrong")
	}

}
