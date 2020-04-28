package connection

import (
	"fmt"
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"strconv"

	"github.com/google/uuid"
)

const (
	// ActionAdmin is an admin type of event
	ActionAdmin = "admin"

	// ActionGame is a game type of event
	ActionGame = "game"

	// ActionMsg is a message type of event
	ActionMsg = "message"

	// ActionPlayerRegister is the register action
	ActionPlayerRegister = "register"

	// ActionPlayerLeave is the leave action
	ActionPlayerLeave = "leave"
)

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
	Action   string
	GameID   int64
	PlayerID int64
	Value    interface{}
}

// ValidateAuthorized ensures game owner is making admin decision
func (e AdminEvent) ValidateAuthorized(game *models.Game) error {
	switch e.Action {
	case ClientChipRequest:
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
func (e AdminEvent) GetChipRequest() ChipRequest {
	value := int64(0)
	switch e.Value.(type) {
	case float64:
		value = int64(e.Value.(float64))
		break
	case int64:
		value = e.Value.(int64)
		break
	case string:
		val, err := strconv.Atoi(e.Value.(string))
		if err != nil {
			fmt.Printf("Error handling chip request value: %s\n", err)
		}
		value = int64(val)
	default:
		value = e.Value.(int64)
	}

	return ChipRequest{uuid.New().String(), e.PlayerID, int64(value)}
}

// MsgEvent represents a message event
type MsgEvent struct {
	Action   string
	GameID   int64
	PlayerID int64
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
