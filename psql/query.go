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

func InsertMatch(tx *sql.Tx, ctx *context.Context, match *elo.Match) error {
	
	insertMatchStmt := `INSERT INTO MATCH 
	(player_a_ID, player_b_ID, player_won_ID, player_a_rating, player_b_rating)
	 VALUES ($1, $2, $3, $4, $5)`
	
	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()
	

	results, err := tx.ExecContext(newCtx, insertMatchStmt, 
		match.PlayerA.Player_ID.String(), match.PlayerB.Player_ID.String(), 
		match.PlayerWon.Player_ID.String(), match.PlayerARating, match.PlayerBRating)
	
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

func GetPlayer(db *sql.DB, ctx *context.Context, name string) (*elo.Player, error) {
	sqlStmt := `Select * FROM player WHERE name=$1`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	row := db.QueryRowContext(newCtx, sqlStmt, name)

	player := &elo.Player{}

	err := row.Scan(player.Player_ID, player.Name, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches)
	
	if err != nil {
		return nil, err	
	}

	return player, nil
}

func GetPlayerWithTX(tx *sql.Tx, ctx *context.Context, name string) (*elo.Player, error) {
	sqlStmt := `Select * FROM player WHERE name=$1`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	row := tx.QueryRowContext(newCtx, sqlStmt, name)

	player := &elo.Player{}

	err := row.Scan(player.Player_ID, player.Name, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches)
	
	if err != nil {
		return nil, err	
	}

	return player, nil
}

func UpdatePlayerWithTx(tx *sql.Tx, ctx *context.Context, player *elo.Player) error {
	sqlStmt :=	`UPDATE player 
	SET elo_rating=$1, wins=$2, losses=$3, draws=$4, total_matches=$5
	WHERE player_ID=$6`

	newCtx, cancel := context.WithTimeout(*ctx, 5*time.Second)
	defer cancel()

	result, err := tx.ExecContext(newCtx, sqlStmt, player.EloRating, player.Wins, player.Losses, player.Draws, player.TotalMatches, player.Player_ID.String())
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
