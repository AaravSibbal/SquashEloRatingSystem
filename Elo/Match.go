package elo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Match struct {
	id *uuid.UUID
	playerA    *Player
	playerB    *Player
	playerWon  *Player
	when       *time.Time
}

type Matches struct {
	matches []*Match
}

func (ms *Matches) New(playerA, playerB, playerWon *Player) *Match {
	time := time.Now()
	id := uuid.New()
	
	m := &Match{
		id: &id,
		playerA: playerA,
		playerB: playerB,
		playerWon: playerWon,
		when: &time,
	}

	playerA.EloRating = GetNewElo(playerA, m)
	playerB.EloRating = GetNewElo(playerB, m)

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
		if match.id == id {
			return idx
		}
	}

	return -1
}