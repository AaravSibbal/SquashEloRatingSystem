package elo

import (
	"fmt"
	"strings"

)

type Player struct {
	Player_ID    string
	Discord_ID   string
	Name         string
	EloRating    int
	Wins         int
	Losses       int
	Draws        int
	TotalMatches int
}

func (player *Player) Equals(p *Player) bool {
	result := strings.Compare(player.Discord_ID, p.Discord_ID)
	return result == 0
}

func (player *Player) UpdatePlayer(match *Match) {
	
	result := player.DidPlayerWin(match)
	player.EloRating = GetNewElo(player, match)
	if result {
		player.Wins++
	} else {
		player.Losses++
	}
	player.TotalMatches++
}

func (player *Player) DidPlayerWin(match *Match) bool {
	return player.Equals(match.PlayerWon)
}

func (player *Player) String() string {
	return fmt.Sprintf("Name: %s\nElo Rating: %d\nWins: %d\nLosses: %d\nTotal Matches: %d\n",
		player.Name, player.EloRating, player.Wins, player.Losses, player.TotalMatches)
}

func (player *Player) StartingEloFromLevel(level string) int {
	if level == ""{
		return 400
	} else if level == "1" {
		return 400
	} else if level == "2" {
		return 500
	} else if level == "3" {
		return 600
	}

	return 400
}

func (player *Player) New(discordID, name, level string){
	player.Name = name
	player.Discord_ID = discordID
	player.EloRating = player.StartingEloFromLevel(level)
	player.Wins = 0
	player.Losses = 0
	player.Draws = 0
	player.TotalMatches = 0
}