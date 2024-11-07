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

// TODO: add beginner, intermidiate, advance option to the addPlayer
func (b *Bot) newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	fmt.Printf("Message Content: %s\n", message.Content)
	switch {
	case strings.HasPrefix(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, b.getHelpMessage())

	case strings.HasPrefix(message.Content, "!ping"):
		discord.ChannelMessageSend(message.ChannelID, b.ping())

	case strings.HasPrefix(message.Content, "!addPlayer"):
		discord.ChannelMessageSend(message.ChannelID, b.addPlayer(message.Mentions))
	case strings.HasPrefix(message.Content, "!addMatch"):
		discord.ChannelMessageSend(message.ChannelID, b.addMatch(message.Mentions, message.Content))
	case strings.HasPrefix(message.Content, "!stat"):
		discord.ChannelMessageSend(message.ChannelID, b.stat(message.Mentions))
	// default:
	// 	discord.ChannelMessageSend(message.ChannelID, "something went wrong")
	}
}
