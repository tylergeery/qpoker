package connection

import (
	"qpoker/cards"
	"qpoker/cards/games/holdem"
)

const (
	// ActionAdmin is an admin type of event
	ActionAdmin = "admin"

	// ActionAdminStart is an admin start of game
	ActionAdminStart = "start"

	// ActionGame is a game type of event
	ActionGame = "game"

	// ActionMsg is a message type of event
	ActionMsg = "message"

	// ActionPlayerRegister is the register action
	ActionPlayerRegister = "register"

	// ActionPlayerLeave is the leave action
	ActionPlayerLeave = "leave"
)

// GameController handles logic for sending/receiving game events
type GameController struct {
	clients []*Client
	manager *holdem.GameManager
}

// GameState controls the game state returned to clients
type GameState struct {
	Manager *holdem.GameManager     `json:"manager"`
	Cards   map[int64][]cards.Card  `json:"cards"`
	Players map[int64]holdem.Player `json:"players"`
}

// NewGameState returns the game state for clients
func NewGameState(manager *holdem.GameManager) GameState {
	return GameState{
		Manager: manager,
		Cards:   manager.GetVisibleCards(),
	}
}

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
	GameID   int64
	PlayerID int64
	Action   holdem.Action
}

// AdminEvent represent an admin action
type AdminEvent struct {
	Action string
	GameID int64
}

// MsgEvent represents a message event
type MsgEvent struct {
}
