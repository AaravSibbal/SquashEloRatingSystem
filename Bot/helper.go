package bot

import (
	"errors"
	"fmt"
	// "fmt"
	"strings"

	elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) GetPlayerWon(playerA *elo.Player, playerB *elo.Player, playerWon string) *elo.Player {
	if playerA.Name == playerWon {
		return playerA
	}

	return playerB
}

func (b *Bot) GetPlayerWonByID(playerWonID string, playerA *discordgo.User, playerB *discordgo.User) *discordgo.User {
	if playerA.ID == playerWonID {
		return playerA
	}

	return playerB
}

func (b *Bot) TrimID(idWithJunk string) string{
	asRunes := []rune(idWithJunk)
	asRunes = asRunes[2 : len(asRunes)-1]

	return string(asRunes)
}

func (b *Bot) UserEqual(user *discordgo.User, id string) bool {
	return user.ID == id
}

func (b *Bot) GetPlayers(users []*discordgo.User, message string)(*discordgo.User, *discordgo.User, *discordgo.User, error) {
	arr := strings.Split(message, " ")
	
	if len(arr) != 4 {
		// fmt.Printf("Arr: %v", arr)
		return nil, nil, nil, fmt.Errorf("Expecting 4 arguments for %d", len(arr))
	}
	if len(users) != 2 {
		return nil, nil, nil, errors.New("playerA, playerB and player won are all different")
	}

	playerAId := b.TrimID(arr[1])
	playerBId := b.TrimID(arr[2])
	playerWonID := b.TrimID(arr[3])

	if playerAId == playerBId {
		return nil, nil, nil, errors.New("playerA and playerB can't be the same")
	}	
	if playerAId != playerWonID && playerBId != playerWonID {
		return nil, nil, nil, errors.New("playerA and playerB don't match playerWon")
	}

	var playerA *discordgo.User
	var playerB *discordgo.User
	var playerWon *discordgo.User

	playerA = users[0]
	playerB = users[1]
	playerWon = b.GetPlayerWonByID(playerWonID, playerA, playerB)

	if b.isBot(playerA) || b.isBot(playerB) || b.isBot(playerWon) {
		return nil, nil, nil, errors.New("can't add a match with a bot")
	}

	return playerA, playerB, playerWon, nil

}

func (b *Bot) isBot(user *discordgo.User) bool {
	return user.Bot
}

func (b *Bot) getLevelFromMessage(message string) string {
	arr := strings.Split(message, " ")
	if len(arr) < 3 {
		return ""
	} 

	return arr[2]
}