# QPoker

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
- Add migration system