CREATE TABLE IF NOT EXISTS player (
    id serial PRIMARY KEY,
    username text NOT NULL,
    email text NOT NULL,
    pw text NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_username UNIQUE(username),
    CONSTRAINT unique_email UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS game (
    id serial PRIMARY KEY,
    name text NOT NULL,
    slug text NOT NULL,
    owner_id integer NOT NULL REFERENCES player(id),
    capacity integer DEFAULT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_slug UNIQUE(slug)
);

CREATE TABLE IF NOT EXISTS game_player_history (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    player_id integer NOT NULL REFERENCES player(id),
    started_at timestamptz DEFAULT NOW(),
    ended_at timestamptz DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS game_table (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    name text NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_table_name UNIQUE(game_id, name)
);

CREATE TABLE IF NOT EXISTS game_table_player (
    id serial PRIMARY KEY,
    game_table_id integer NOT NULL REFERENCES game_table(id),
    seat_index smallint NOT NULL,
    player_id integer NOT NULL REFERENCES player(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_table_seat UNIQUE(game_table_id, seat_index),
    CONSTRAINT game_table_player_id UNIQUE(game_table_id, player_id)
);

CREATE TABLE IF NOT EXISTS game_table_player_history (
    id serial PRIMARY KEY,
    game_table_id integer NOT NULL REFERENCES game(id),
    player_id integer NOT NULL REFERENCES player(id),
    started_at timestamptz DEFAULT NOW(),
    ended_at timestamptz DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS game_table_hand (
    id serial PRIMARY KEY,
    game_table_id integer NOT NULL REFERENCES game_table(id),
    hand_index integer NOT NULL,
    board text[],
    state text NOT NULL,
    winner integer NOT NULL REFERENCES player(id),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_table_hand_index UNIQUE(game_table_id, hand_index)
);

CREATE TABLE IF NOT EXISTS game_table_player_hand (
    id serial PRIMARY KEY,
    game_table_hand_id integer NOT NULL REFERENCES game_table_hand(id),
    player_id integer NOT NULL REFERENCES game_table_player_history(id),
    cards text[],
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_table_hand_player_constraint UNIQUE(game_table_hand_id, player_id)
);
