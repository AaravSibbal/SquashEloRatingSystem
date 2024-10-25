package bot

import elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"

func (b *Bot) GetPlayerWon(playerA *elo.Player, playerB *elo.Player, playerWon string) *elo.Player {
	if playerA.Name == playerWon {
		return playerA
	}

	return playerB
}