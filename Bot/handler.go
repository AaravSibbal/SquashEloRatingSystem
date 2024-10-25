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
	if len(users) != 1 {
		return fmt.Sprintf("Error: expected 1 users got %d", len(users))
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

func (b *Bot) addMatch(users []*discordgo.User) string {
	if len(users) != 3 {
		return fmt.Sprintf("Error: expected 3 mentions got %d", len(users))
	}

	tx, err := b.Db.BeginTx(*b.Ctx, nil)

	if err != nil {
		return "there was some trouble with the db try again later"
	}

	playerAName := users[0].GlobalName
	playerBName := users[1].GlobalName
	playerWonName := users[2].GlobalName

	if playerAName != playerWonName && playerBName != playerWonName {
		return "Error: player won name does not match playerA or playerB"

	}

	playerA, err := psql.GetPlayerWithTX(tx, b.Ctx, playerAName)
	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerAName)
	}

	playerB, err := psql.GetPlayerWithTX(tx, b.Ctx, playerBName)
	if err != nil {
		tx.Rollback()
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerAName)
	}

	playerWon := b.GetPlayerWon(playerA, playerB, playerWonName)

	match := elo.Matches.New(elo.Matches{}, playerA, playerB, playerWon)

	err = psql.InsertMatch(tx, b.Ctx, match)
	if err != nil {
		tx.Rollback()
		return "Couldn't add match to the db, my bad lol..."
	}		
	tx.Commit()
	return "Added the match successfully"
}