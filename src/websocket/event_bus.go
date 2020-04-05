package main

import (
	"fmt"
)

const (
	actionPlayerRegister = "register"
	actionPlayerLeave    = "leave"

	actionGameBet   = "bet"
	actionGameCheck = "check"
	acitonGameFold  = "fold"
)

// PlayerEvent represents a player connection action
type PlayerEvent struct {
	Client *Client
	Action string
}

// GameEvent represents a player gameplay action
type GameEvent struct {
	Client *Client
	Action string
	Amount int64
}

// EventBus manages all server event action
type EventBus struct {
	clients       map[int64]*Client
	PlayerChannel chan PlayerEvent
	GameChannel   chan GameEvent
}

// NewEventBus returns a new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		clients: map[int64]*Client{},
	}
}

func (e *EventBus) setClient(client *Client) {
	e.clients[client.PlayerID] = client
}

func (e *EventBus) removeClient(playerID int64) {
	delete(e.clients, playerID)
}

// ListenForEvents starts the event bus waiting for channel events
func (e *EventBus) ListenForEvents() {
	for {
		select {
		case playerEvent := <-e.PlayerChannel:
			fmt.Printf("PlayerAction: (%d %s)\n", playerEvent.Client.PlayerID, playerEvent.Action)
			if playerEvent.Action == actionPlayerRegister {
				e.setClient(playerEvent.Client)
			}
			if playerEvent.Action == actionPlayerLeave {
				e.removeClient(playerEvent.Client.PlayerID)
			}
		case gameEvent := <-e.GameChannel:
			fmt.Printf("GameAction: (%d %s)\n", gameEvent.Client.PlayerID, gameEvent.Action)
			// TODO: validate game action, update game state
		}
	}
}
