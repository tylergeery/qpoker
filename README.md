# QPoker

A live video multi-player Texas HoldEm Poker game.

Host your own poker game with your friends.

Try it out [here]()

## Supported Card Games
- Holdem (Poker)
- Hearts

## Hosting / Deployment Support
QCards supports Docker runtimes deployable with either

- Ansible
- k8s

## Architecture
QCards compromises 2 publicly exposed Go services, dependent upon shared Postgres and Redis databases:
- http server
  - handles REST-like API support for games/users
- websocket server
  - handles real-time event communication:
    - game events
    - video
    - messages

## REST API

| Method | URL                             | Description                       
|--------|---------------------------------|-----------------------------------
| GET    | /api/v1/players/:id             | Get player
| POST   | /api/v1/players                 | Create player
| POST   | /api/v1/players/login           | Player login
| PUT    | /api/v1/players/:id             | Update player
| GET    | /api/v1/games/:game_id          | Get game
| POST   | /api/v1/games                   | Create a game
| PUT    | /api/v1/games/:game_id          | Update a game
| GET    | /api/v1/games/:game_id/history  | Get game history

## Development

For available commands, run
```bash
make
```

## Run Test Suite
```bash
make test
```

## Contributions
Please. No guidelines yet.

## Special thanks to (built on):
- [material.css]()
- [card svgs]()
