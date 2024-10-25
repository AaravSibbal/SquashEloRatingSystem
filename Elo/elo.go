package elo

import "math"

func GetNewElo(p *Player, match *Match) int {
	var opponent *Player
	var player *Player

	if p.Equals(match.PlayerA) {
		player = match.PlayerA
		opponent = match.PlayerB
	} else {
		player = match.PlayerB
		opponent = match.PlayerA
	}

	expectedScore := CalculateExpectedScore(player.EloRating, opponent.EloRating)

	if player.Equals(match.PlayerWon){
		return CalculateElo(player.EloRating, expectedScore, true)
	} else {
		return CalculateElo(player.EloRating, expectedScore, false)
	}
}

func CalculateExpectedScore(playerRating int, opponentRating int) float64 {
	return 1 / (1 + math.Pow(10,((float64(opponentRating)-float64(playerRating))/400)))
}

func CalculateElo(eloRating int, expectedScore float64, didWin bool) int {
	var actualScore float64

	if didWin {
		actualScore = 1 
	} else {
		actualScore = 0
	}
	
	KFactor := CalculateKFactor(eloRating)
	
	newElo := eloRating + int(math.Ceil((KFactor * (actualScore - expectedScore))))
	
	return minEloCheck(newElo)
}

func CalculateKFactor(eloRating int) float64 {
	if(eloRating < 500){
		return 32
	} else if eloRating < 700 {
		return 24
	} else if eloRating < 1000 {
		return 20
	} else if eloRating < 1300 {
		return 16
	} else {
		return 16
	}
	
}

func minEloCheck(newElo int) int {
	if newElo < 200 {
		return 200
	}

	return newElo
}