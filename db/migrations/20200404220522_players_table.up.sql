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
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_slug UNIQUE(slug)
);

CREATE TABLE IF NOT EXISTS game_options (
    id serial PRIMARY KEY,
    game_id integer NOT NULL REFERENCES game(id),
    options jsonb NOT NULL,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    CONSTRAINT unique_game_options_game_id UNIQUE(game_id)
);
