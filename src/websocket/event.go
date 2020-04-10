package main

// Event is the event broadcasted to all clients
type Event struct {
	Type  string `json:"type"`
	State string `json:"state"`
}

// PlayerEvent represents a player connection action
type PlayerEvent struct {
	Client *Client
	Action string
}

// AuthEvent is a client initiated event to verify game
type AuthEvent struct {
	Token  string `json:"token"`
	GameID int64  `json:"game_id"`
}

// GameEvent represents a player gameplay action
type GameEvent struct {
	PlayerID int64  `json:"-"`
	Action   string `json:"action"`
	Amount   int64  `json:"amount"`
}

// AdminEvent represent an admin gameplay action
type AdminEvent struct {
	Action string
}
