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
      - Hosting options
      - Card games available
    - Table page
      - Table rendering
  - JS
    - Table page
      - Forced Login
      - Game flow
      - Game History
    - Add tests
  - Assets
    - Card assets
    - Table asset
- Admin 
  - Seat assignments
  - Assigning Chips
  - Time Between hands
- EventTracking
  - action events
  - admin events
  - chat events
- Write GameManager state to DB (in case of disconnect)
- Deploy somewhere
- Autovaccuum games
- Video Chat
- Add use options
  - Charge for hosting
  - Free hosting w/ ads
  - Host your own