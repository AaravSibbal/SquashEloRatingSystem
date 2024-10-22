CREATE TABLE IF NOT EXISTS player(
    player_ID UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE, 
    elo_rating SMALLINT NOT NULL, 
    wins INT NOT NULL,
    losses INT NOT NULL,
    draws INT NOT NULL,
    total_matches INT NOT NULL
);

CREATE TABLE IF NOT EXISTS match(
    match_id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    player_a_ID UUID NOT NULL,
    player_b_ID UUID NOT NULL,
    player_won_ID UUID NOT NULL,
    match_time timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_player_a_ID
        FOREIGN KEY(player_a_ID)
            REFERENCES player(player_ID),
    CONSTRAINT fk_player_b_ID
        FOREIGN KEY(player_b_ID)
            REFERENCES player(player_ID),
    CONSTRAINT fk_player_won_ID
        FOREIGN KEY(player_won_ID)
            REFERENCES player(player_ID)
);



	-- id *uuid.UUID
	-- playerA    *Player
	-- playerB    *Player
	-- playerWon  *Player
	-- when       *time.Time
	
    -- Name         string
	-- EloRating    int
	-- Wins         int
	-- Losses       int
	-- Draws        int
	-- TotalMatches int


CREATE TABLE IF NOT EXISTS blogs(
    blog_ID UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    title VARCHAR(1023) NOT NULL UNIQUE,
    fileHandle VARCHAR(2047) NOT NULL UNIQUE,
    created timestamptz NOT NULL DEFAULT now()
);