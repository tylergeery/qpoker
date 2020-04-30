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
  - Landing page
    - Join game by code
    - Features
      - Video
      - Card games
        - holdem
        - hearts
        - gang gin
      - self-hosting
    - Hosting options
      - Charge for hosting
      - Free hosting w/ ads
      - Host your own
  - JS
    - Table page
      - Game History
      - Remove nav choices
    - Add tests
- Admin Options
  - Time Between hands
  - Time limit on user choice
- Deploy somewhere
- Autovaccuum games
  - Rules around when games should be removed from memory
  - How to handle client disconnects
    - immediately, delay for reconnect?
- Video Chat