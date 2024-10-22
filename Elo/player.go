package elo

import (
	"strings"

	"github.com/google/uuid"
)

type Player struct {
	Player_ID 	*uuid.UUID
	Name         string
	EloRating    int
	Wins         int
	Losses       int
	Draws        int
	TotalMatches int
}

func (player *Player) Equals(p *Player) bool {
	result := strings.Compare(player.Name, p.Name)
	return result == 0
}

type Players struct {
	players map[string]*Player
}

func (pls *Players) New(name string) *Player {
	p := &Player{
		Name:         name,
		EloRating:    400,
		Wins:         0,
		Losses:       0,
		Draws:        0,
		TotalMatches: 0,
	}
	return p
}

func (pls *Players) AddPlayer(p *Player) bool {
	_, exists := pls.players[p.Name]
	if exists {
		return false
	}

	pls.players[p.Name] = p
	return true
}

func (pls *Players) RemovePlayer(name string) bool {
	_, exists := pls.players[name]
	if !exists {
		return false
	}

	delete(pls.players, name)
	return true
}

func (pls *Players) GetPlayer(name string) (*Player, bool) {
	player, exists := pls.players[name]
	if !exists {
		return nil, false
	}
	return player, true
}
