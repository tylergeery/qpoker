# QPoker

A live video multi-player Texas HoldEm Poker game.

Host your own poker game with your friends.


## API

| Method | URL                             | Description                       
|--------|---------------------------------|-----------------------------------
| GET    | /api/v1/users/:user_id                         | Get user
| GET    | /api/v1/users/:user_id/games                   | Get a users' games
| POST   | /api/v1/users                                  | Create user
| GET    | /api/v1/games/:game_id                         | Get game
| POST   | /api/v1/games                                  | Create a game
| GET    | /api/v1/games/:game_id/tables                  | Get all tables in a given game
| GET    | /api/v1/games/:game_id/tables/:table_id        | Get table in game


## Development

For available commands, run
```bash
make
```

## Run Test Suite
```bash
make test
```

## TODO
- Create all API Endpoints
  - Add tests
- Add tests for websocket server
- Add static pages
- Add JS client
- Features
  - Holdem scoring
  - Assigning chips
  - Assets
    - Card assets
    - Table Asset
  - Video Chat