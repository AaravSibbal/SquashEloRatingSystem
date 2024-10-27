package bot

import (
	"fmt"

	elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"
	"github.com/AaravSibbal/SqashEloRatingSystem/psql"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getHelpMessage() string {
	fmt.Println("got the help message")
	return `help is not implemented yet`
}

func (b *Bot) ping() string {
	return "pong"
}

func (b *Bot) addPlayer(users []*discordgo.User) string {

	if len(users) != 2 {
		return fmt.Sprintf("Error: expected 2 users got %d", len(users))
	}

	user := users[0]
	if user.Bot {
		return "Error: Player can't be a bot"
	}

	fmt.Printf("adding player: %v", user.GlobalName) 
	player := elo.Players.New(elo.Players{}, user.GlobalName)
	err := psql.InsertPlayer(b.Db, b.Ctx, player)

	if err != nil {
		return "there was an error adding the player, " + err.Error()
	}

	return fmt.Sprintf("Player: %s, was added successfully", user.GlobalName)
}

func (b *Bot) addMatch(users []*discordgo.User, message string) string {
	playerAUser, playerBUser, playerWonUser, err := b.GetPlayers(users, message)
	if err != nil {
		return err.Error()
	}

	tx, err := b.Db.BeginTx(*b.Ctx, nil)

	if err != nil {
		return "there was some trouble with the db try again later"
	}

	playerA, err := psql.GetPlayerWithTX(tx, b.Ctx, playerAUser.GlobalName)
	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerAUser.GlobalName)
	}

	playerB, err := psql.GetPlayerWithTX(tx, b.Ctx, playerBUser.GlobalName)
	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerBUser.GlobalName)
	}

	playerWon := b.GetPlayerWon(playerA, playerB, playerWonUser.GlobalName)

	match := elo.Matches.New(elo.Matches{}, playerA, playerB, playerWon)

	err = psql.InsertMatch(tx, b.Ctx, match)
	if err != nil {
		tx.Rollback()
		return "Couldn't add match to the db, my bad lol..."
	}		
	tx.Commit()
	return "Added the match successfully"
}

func (b *Bot) stat(users []*discordgo.User) string {
	if len(users) < 1 {
		return "need at least 1 user in order get the stat"
	}

	resultStr := ""

	for _, user := range users {
		resultStr += b.getPlayerStat(user)
	}

	return resultStr 
}

func (b *Bot) getPlayerStat(user *discordgo.User) string {
	
}