package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"
)

func InsertPlayer(db *sql.DB, ctx *context.Context, player *elo.Player) error {
	sqlStmt := `INSERT INTO player 
	(name, elo_rating, wins, losses, draws, total_matches, discord_ID) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(newCtx, sqlStmt, player.Name, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches, player.Discord_ID)
	if err == context.DeadlineExceeded {
		fmt.Printf("The query took too long for InsertPlayer, %v\n", err)
		return err
	} else if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil{
		return err
	}
	if rows != 1 {
		return fmt.Errorf("expected row affected to be 1 got %d", rows)
	}


	return nil
}

func InsertMatch(tx *sql.Tx, ctx *context.Context, match *elo.Match) error {
	
	insertMatchStmt := `INSERT INTO MATCH 
	(player_a_ID, player_b_ID, player_won_ID, player_a_rating, player_b_rating,
	player_a_name, player_b_name)
	 VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()
	

	results, err := tx.ExecContext(newCtx, insertMatchStmt, 
		match.PlayerA.Player_ID, match.PlayerB.Player_ID, 
		match.PlayerWon.Player_ID, match.PlayerARating, match.PlayerBRating,
		match.PlayerAName, match.PlayerBName)
	
		if err == context.DeadlineExceeded {
		return err
	} else if err != nil {
		return err
	}

	row, err := results.RowsAffected()
	if err == sql.ErrNoRows {
		return err
	} else if row != 1 {
		return fmt.Errorf("expected 1 row to be affected but rows affected were %v", row)
	}

	return nil
}

/*
uses the discord id of user to fetch the user return the Player obj if 
we get a valid result return an error otherwise
*/
func GetPlayer(db *sql.DB, ctx *context.Context, discordID string) (*elo.Player, error) {
	sqlStmt := `Select * FROM player WHERE discord_ID=$1`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(newCtx, sqlStmt, discordID)

	player := &elo.Player{}

	err := row.Scan(&player.Player_ID, &player.Name, &player.EloRating, &player.Wins, &player.Losses, &player.Draws, &player.TotalMatches, &player.Discord_ID)
	
	if err != nil {
		return nil, err	
	}

	return player, nil
}

func UpdatePlayerWithTx(tx *sql.Tx, ctx *context.Context, player *elo.Player) error {
	sqlStmt :=	`UPDATE player 
	SET elo_rating=$1, wins=$2, losses=$3, draws=$4, total_matches=$5
	WHERE discord_ID=$6`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	result, err := tx.ExecContext(newCtx, sqlStmt, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches, player.Discord_ID) 
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("more rows were affected than expected")
	}

	return nil
}
