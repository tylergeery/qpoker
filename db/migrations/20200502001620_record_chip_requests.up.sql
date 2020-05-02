CREATE TABLE IF NOT EXISTS game_chip_requests (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    player_id integer NOT NULL REFERENCES player(id),
    amount integer NOT NULL,
    status text NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

ALTER TABLE game ADD COLUMN status text NOT NULL;
