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
	Db       *sql.DB
	Ctx      *context.Context
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

func standardizeSpaces(s string) string {
    return strings.Join(strings.Fields(s), " ")
}

// TODO: add beginner, intermidiate, advance option to the addPlayer
// TODO: make it so you are actually storing the id of the player
func (b *Bot) newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	message.Content = strings.ToLower(message.Content)
	message.Content = standardizeSpaces(message.Content)


	fmt.Printf("Message Content: %s\n", message.Content)
	switch {
	case strings.HasPrefix(message.Content, "!help"):
		discord.ChannelMessageSend(message.ChannelID, b.getHelpMessage())

	case strings.HasPrefix(message.Content, "!ping"):
		discord.ChannelMessageSend(message.ChannelID, b.ping())

	case strings.HasPrefix(message.Content, "!addplayer"):
		discord.ChannelMessageSend(message.ChannelID, b.addPlayer(message.Mentions, message.Content))
	case strings.HasPrefix(message.Content, "!addmatch"):
		discord.ChannelMessageSend(message.ChannelID, b.addMatch(message.Mentions, message.Content))
	case strings.HasPrefix(message.Content, "!stat"):
		discord.ChannelMessageSend(message.ChannelID, b.stat(message.Mentions))
		// default:
		// 	discord.ChannelMessageSend(message.ChannelID, "something went wrong")
	}
}
