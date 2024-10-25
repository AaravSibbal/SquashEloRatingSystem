package elo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	Id *uuid.UUID
	PlayerA    *Player
	PlayerB    *Player
	PlayerWon  *Player
	PlayerARating int
	PlayerBRating int
	When       *time.Time
}

type Matches struct {
	matches []*Match
}

func (ms Matches) New(playerA, playerB, playerWon *Player) *Match {
	time := time.Now()
	id := uuid.New()
	
	m := &Match{
		Id: &id,
		PlayerA: playerA,
		PlayerB: playerB,
		PlayerWon: playerWon,
		PlayerARating: playerA.EloRating,
		PlayerBRating: playerB.EloRating,
		When: &time,
	}

	
	// this step is important because otherwise elo would change while it is getting calculated
	playerANewRating := GetNewElo(playerA, m)
	playerBNewRating := GetNewElo(playerB, m)

	playerA.EloRating = playerANewRating
	playerB.EloRating = playerBNewRating

	playerA.UpdatePlayer(m)
	playerB.UpdatePlayer(m)

	return m
}

func (ms *Matches) AddMatch(match *Match) bool {
	if match == nil {
		return false
	}

	ms.matches = append(ms.matches, match)
	return true
}

func (ms *Matches) RemoveMatch(id *uuid.UUID) (bool, error){
	matchIdx := ms.GetMatchIdx(id)
	if matchIdx == -1 {	
		return false, errors.New("couldn't find the match")
	}
	if matchIdx < 0 || matchIdx > len(ms.matches){
		return false, errors.New("index out of bounds for matches")
	}
	
	ms.matches = append(ms.matches[:matchIdx], ms.matches[matchIdx+1:]...)
	return true, nil
}


func (ms *Matches) GetMatchIdx(id *uuid.UUID) int {
	for idx,match := range ms.matches {
		if match.Id == id {
			return idx
		}
	}

	return -1
}