CREATE TABLE IF NOT EXISTS game_hand (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    board text[],
    state text NOT NULL,
    payouts json DEFAULT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game_player_hand (
    id serial PRIMARY KEY,
    game_hand_id integer NOT NULL REFERENCES game_hand(id),
    player_id integer NOT NULL REFERENCES player(id),
    cards text[],
    cards_visible boolean NOT NULL,
    starting_stack integer NOT NULL,
    ending_stack integer NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_game_player_hand_constraint UNIQUE(game_hand_id, player_id)
);
