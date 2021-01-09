package events

import (
	"fmt"
	"qpoker/models"
	"qpoker/websocket/connection"
	"qpoker/websocket/utils"
)

const (
	// ActionAdmin is an admin type of event
	ActionAdmin = "admin"

	// ActionGame is a game type of event
	ActionGame = "game"

	// ActionMsg is a message type of event
	ActionMsg = "message"

	// ActionVideo is a video type of event
	ActionVideo = "video"

	// ActionPlayerRegister is the register action
	ActionPlayerRegister = "register"

	// ActionPlayerLeave is the leave action
	ActionPlayerLeave = "leave"
)

// Chat message
type Chat struct {
	PlayerID       int64  `json:"player_id"`
	PlayerUsername string `json:"player_username"`
	Message        string `json:"message"`
}

// PlayerEvent represents a player connection action
type PlayerEvent struct {
	Client *connection.Client
	Action string
}

// GameEvent represents a player gameplay action
type GameEvent struct {
	GameID   int64
	PlayerID int64
	Action   interface{}
}

// AdminEvent represent an admin action
type AdminEvent struct {
	Action   string
	GameID   int64
	PlayerID int64
	Value    interface{}
}

// ValidateAuthorized ensures game owner is making admin decision
func (e AdminEvent) ValidateAuthorized(game *models.Game) error {
	switch e.Action {
	case connection.ClientChipRequest:
		return nil
	default:
		if e.PlayerID != game.OwnerID {
			fmt.Println("Invalid authorization")
			return fmt.Errorf("Invalid authorization for (%d), expected (%d)", e.PlayerID, game.OwnerID)
		}
	}

	return nil
}

// GetChipRequest gets chip request from event
func (e AdminEvent) GetChipRequest() *models.GameChipRequest {
	return &models.GameChipRequest{
		GameID:   e.GameID,
		PlayerID: e.PlayerID,
		Amount:   utils.InterfaceInt64(e.Value),
		Status:   models.GameChipRequestStatusInit,
	}
}

// MsgEvent represents a message event
type MsgEvent struct {
	Action   string
	GameID   int64
	PlayerID int64
	Value    string
	Username string
}

// GetChat turns MsgEvent into Chat
func (e MsgEvent) GetChat() Chat {
	return Chat{
		PlayerID:       e.PlayerID,
		Message:        e.Value,
		PlayerUsername: e.Username,
	}
}

// BroadcastEvent is the event broadcasted to all clients
type BroadcastEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewBroadcastEvent creates a new BroadcastEvent
func NewBroadcastEvent(eventType string, data interface{}) BroadcastEvent {
	return BroadcastEvent{eventType, data}
}

// VideoEvent represents a message event
type VideoEvent struct {
	Type         string      `json:"type"`
	GameID       int64       `json:"game_id"`
	FromPlayerID int64       `json:"from_player_id"`
	ToPlayerID   int64       `json:"to_player_id"`
	Offer        interface{} `json:"offer"`
	Candidate    interface{} `json:"candidate"`
}
