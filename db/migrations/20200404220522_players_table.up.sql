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
