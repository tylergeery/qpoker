CREATE TABLE IF NOT EXISTS game_hand (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    board text[] DEFAULT NULL,
    payouts jsonb DEFAULT NULL,
    bets jsonb DEFAULT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game_player_hand (
    id serial PRIMARY KEY,
    game_hand_id integer NOT NULL REFERENCES game_hand(id),
    player_id integer NOT NULL REFERENCES player(id),
    cards text[] DEFAULT NULL,
    starting integer NOT NULL,
    ending integer DEFAULT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_game_player_hand_constraint UNIQUE(game_hand_id, player_id)
);
