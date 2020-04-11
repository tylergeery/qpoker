package main

import (
	// "encoding/json"
	"fmt"

	"qpoker/cards/games/holdem"
	"qpoker/models"
)

const (
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

// EventBus manages all server event action
type EventBus struct {
	games         map[int64]*GameController
	PlayerChannel chan PlayerEvent
	GameChannel   chan GameEvent
	AdminChannel  chan AdminEvent
}

// NewEventBus returns a new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		games:         map[int64]*GameController{},
		PlayerChannel: make(chan PlayerEvent),
		GameChannel:   make(chan GameEvent),
	}
}

func (e *EventBus) reloadGameState(client *Client) (*holdem.GameManager, error) {
	// TODO: reload game state if not present
	players := []*holdem.Player{&holdem.Player{ID: client.PlayerID}}
	options := models.GameOptions{Capacity: 12, BigBlind: 50}
	return holdem.NewGameManager(players, options)
}

func (e *EventBus) setClient(client *Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		manager, err := e.reloadGameState(client)
		if err != nil {
			fmt.Errorf("error creating GameManager: %s", err)
			return
		}

		controller = &GameController{
			clients: []*Client{},
			manager: manager,
		}
		e.games[client.GameID] = controller
	}

	_ = controller.manager.State.Table.AddPlayer(&holdem.Player{ID: client.PlayerID})
	controller.clients = append(controller.clients, client)
}

func (e *EventBus) removeClient(client *Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		return
	}

	for i := range controller.clients {
		if controller.clients[i] == client {
			controller.clients = append(controller.clients[:i], controller.clients[i+1:]...)
		}
	}
}

func (e *EventBus) broadcast(gameID int64) {

}

// ListenForEvents starts the event bus waiting for channel events
func (e *EventBus) ListenForEvents() {
	for {
		select {
		case playerEvent := <-e.PlayerChannel:
			fmt.Printf("PlayerAction: (%d %s)\n", playerEvent.Client.PlayerID, playerEvent.Action)
			if playerEvent.Action == ActionPlayerRegister {
				e.setClient(playerEvent.Client)
			}
			if playerEvent.Action == ActionPlayerLeave {
				e.removeClient(playerEvent.Client)
			}
		case adminEvent := <-e.AdminChannel:
			fmt.Printf("AdminAction: (%s)", adminEvent.Action)
		case gameEvent := <-e.GameChannel:
			fmt.Printf("GameAction: (%s)\n", gameEvent.Action)
			// TODO: validate game action, update game state
		}
	}
}
