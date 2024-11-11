package bot

import (
	"fmt"

	elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"
	"github.com/AaravSibbal/SqashEloRatingSystem/psql"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) getHelpMessage() string {
	return `Commands: 

- !addPlayer @player 
(always mention the player as it will ensure that the name is always correct)

- !addMatch @playerA @playerB @playerWon
(make sure that player won either same as PlayerA or PlayerB)

- !help
(lists the commands available)

- !ping
(resonds back with a pong to check if the bot is still running)

!stat @player1
(responds with player(s) elo, wins, losses, total matches)

!stat @player1 @player2 ...
(you can add as many players you wont that are in the server)`

}

func (b *Bot) ping() string {
	return "pong"
}

func (b *Bot) addPlayer(users []*discordgo.User, message string) string {

	if len(users) != 1 {
		return fmt.Sprintf("Error: expected 1 users got %d", len(users))
	}

	user := users[0]
	if user.Bot {
		return "Error: Player can't be a bot"
	}

	fmt.Printf("adding player: %s, username: %s", user.ID, user.Username)

	player := &elo.Player{}
	player.New(user.ID, user.Username, b.getLevelFromMessage(message))
	fmt.Printf("Player Name: %s\n", player.Name)
	err := psql.InsertPlayer(b.Db, b.Ctx, player)

	if err != nil {
		return "there was an error adding the player, " + err.Error()
	}

	return fmt.Sprintf("Player: %s, was added successfully", user.Username)
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

	playerA, err := psql.GetPlayer(b.Db, b.Ctx, playerAUser.ID)
	if err != nil {
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerAUser.ID)
	}

	playerB, err := psql.GetPlayer(b.Db, b.Ctx, playerBUser.ID)
	if err != nil {
		return fmt.Sprintf("Player: %s, not found add them to the db first", playerBUser.ID)
	}

	playerWon := b.GetPlayerWon(playerA, playerB, playerWonUser.Username)

	match := &elo.Match{}
	match.New(playerA, playerB, playerWon)
	playerA.UpdatePlayer(match)
	playerB.UpdatePlayer(match)

	err = psql.InsertMatch(tx, b.Ctx, match)
	if err != nil {
		tx.Rollback()
		fmt.Printf("error: %v", err)
		return "Couldn't add match to the db, my bad lol..."
	}
	err = psql.UpdatePlayerWithTx(tx, b.Ctx, playerA)
	if err != nil {
		tx.Rollback()
		fmt.Printf("error: %v", err)
		return fmt.Sprintf("Error: Couldn't update %s, match was not added", playerA.Name)
	}
	err = psql.UpdatePlayerWithTx(tx, b.Ctx, playerB)
	if err != nil {
		tx.Rollback()
		fmt.Printf("error: %v", err)
		return fmt.Sprintf("Error: Couldn't update %s, match was not added", playerB.Name)
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
	if user == nil {
		return "user is not defined\n"
	}

	player, err := psql.GetPlayer(b.Db, b.Ctx, user.ID)
	if err != nil {
		return fmt.Sprintf("There was some trouble getting %s from the db", player.Name)
	}

	return player.String()
}
