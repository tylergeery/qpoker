CREATE TABLE IF NOT EXISTS game_type (
    id serial PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    is_active boolean NOT NULL DEFAULT '1',
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_type_unique_key UNIQUE(key)
);

ALTER TABLE game ADD COLUMN game_type_id integer NOT NULL REFERENCES game_type(id);

CREATE TABLE IF NOT EXISTS game_option (
    id serial PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game_type_game_option (
    id serial PRIMARY KEY,
    game_type_id integer NOT NULL REFERENCES game_type(id),
    game_option_id integer NOT NULL REFERENCES game_option(id),
    is_active boolean NOT NULL DEFAULT '1',
    default_value VARCHAR(255) NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS game_game_option_value (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    game_type_game_option_id integer NOT NULL REFERENCES game_type_game_option(id),
    value VARCHAR(255) NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT game_game_option_value_unique_game_option UNIQUE(game_id, game_type_game_option_id)
);
