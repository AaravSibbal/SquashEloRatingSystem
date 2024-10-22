package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	elo "github.com/AaravSibbal/SqashEloRatingSystem/Elo"
	"github.com/google/uuid"
)

func InsertPlayer(db *sql.DB, ctx *context.Context, player *elo.Player) error {
	sqlStmt := `INSERT INTO player 
	(name, elo_rating, wins, losses, draws, total_matches) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	result, err := db.ExecContext(newCtx, sqlStmt, player.Name, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches)
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

func InsertMatch(db *sql.DB, ctx *context.Context, match *elo.Match) error {
	
	playerNameQueryStmt := `SELECT player_ID from player where name = ?`
	insertMatchStmt := `INSERT INTO MATCH 
	(player_a_ID, player_b_ID, player_won_ID)
	 VALUES ($1, $2, $3)`
	
	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()
	
	tx, err := db.BeginTx(newCtx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var playerAID uuid.UUID
	var playerBID uuid.UUID
	var playerWonID uuid.UUID

	// making sure that player A and B are in the db and getting their uuid
	playerARow := tx.QueryRowContext(newCtx, playerNameQueryStmt, match.PlayerA.Name)
	playerBRow := tx.QueryRowContext(newCtx, playerNameQueryStmt, match.PlayerB.Name)

	err = playerARow.Scan(playerAID)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	err = playerBRow.Scan(playerBID)
	if err == sql.ErrNoRows {
		return err
	} else if err != nil {
		return err
	}

	fmt.Printf("PlayerA UUID: %v,\nPlayerB UUID: %v", playerAID, playerBID)

	if match.PlayerWon.Equals(match.PlayerA){
		playerWonID = playerAID
	} else{
		playerWonID = playerBID
	}

	results, err := db.ExecContext(newCtx, insertMatchStmt, playerAID, playerBID, playerWonID)
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
