# QPoker

A live video multi-player Texas HoldEm Poker game.

Host your own poker game with your friends.


## API

| Method | URL                             | Description                       
|--------|---------------------------------|-----------------------------------
| GET    | /api/v1/players/:id             | Get player
| POST   | /api/v1/players                 | Create player
| POST   | /api/v1/players/login           | Player login
| PUT    | /api/v1/players/:id             | Update player
| GET    | /api/v1/games/:game_id          | Get game
| POST   | /api/v1/games                   | Create a game
| PUT    | /api/v1/games/:game_id          | Update a game


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
- Add tests for websocket server
- Client
  - Add static pages
    - Landing page
      - Join game by code
      - Hosting options
        - Charge for hosting
        - Free hosting w/ ads
        - Host your own
      - Card games available
        - Holdem
        - Gang gin
        - Hearts
        - TODO
  - JS
    - Table page
      - Game flow
      - Game History
      - Remove nav choices
    - Add tests
- Admin Options
  - Assigning Chips
  - Time Between hands
  - Selecting seats
- Deploy somewhere
- Autovaccuum games
  - Rules around when games should be removed from memory
  - How to handle client disconnects
    - immediately, delay for reconnect?
- Video Chat