package main

import (
	"encoding/json"
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
	game, err := models.GetGameBy("id", client.GameID)
	if err != nil {
		return nil, err
	}

	// TODO: Check if game is complete
	// TODO: convert model players to game players
	players := []*holdem.Player{&holdem.Player{ID: client.PlayerID}}

	return holdem.NewGameManager(players, game.Options)
}

// SetClient adds client to EventBus
func (e *EventBus) SetClient(client *Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		manager, err := e.reloadGameState(client)
		if err != nil {
			fmt.Errorf("error creating GameManager: %s", err)
			return
		}

		controller = &GameController{[]*Client{}, manager}
		e.games[client.GameID] = controller
	}

	_ = controller.manager.State.Table.AddPlayer(&holdem.Player{ID: client.PlayerID})
	controller.clients = append(controller.clients, client)

	e.BroadcastState(client.GameID)
}

// RemoveClient removes a client from EventBus
func (e *EventBus) RemoveClient(client *Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		return
	}

	for i := range controller.clients {
		if controller.clients[i] == client {
			controller.clients = append(controller.clients[:i], controller.clients[i+1:]...)
		}
	}

	e.BroadcastState(client.GameID)
}

// BroadcastState sends game state to all clients
func (e *EventBus) BroadcastState(gameID int64) {
	controller, ok := e.games[gameID]
	if !ok {
		return
	}

	// TODO: include relevant cards
	state, err := json.Marshal(controller.manager)
	if err != nil {
		fmt.Printf("Error broadcasting game state: %s\n", err)
		return
	}

	for i := range controller.clients {
		controller.clients[i].SendMessage(state)
	}
}

// PerformGameAction sends game state to all clients
func (e *EventBus) PerformGameAction(gameEvent GameEvent) {
	controller, ok := e.games[gameEvent.GameID]
	if !ok {
		return
	}

	complete, err := controller.manager.PlayerAction(gameEvent.PlayerID, gameEvent.Action)
	if err != nil {
		fmt.Printf("Error performing gameEvent: %+v, err: %s\n", gameEvent, err)
		return
	}

	e.BroadcastState(gameEvent.GameID)

	if complete {
		// Start timer for next hand
	}
}

// ListenForEvents starts the event bus waiting for channel events
func (e *EventBus) ListenForEvents() {
	for {
		select {
		case playerEvent := <-e.PlayerChannel:
			fmt.Printf("PlayerAction: (%d %s)\n", playerEvent.Client.PlayerID, playerEvent.Action)
			playerEventMap := map[string]func(*Client){
				ActionPlayerRegister: e.SetClient,
				ActionPlayerLeave:    e.RemoveClient,
			}
			playerEventMap[playerEvent.Action](playerEvent.Client)
		case adminEvent := <-e.AdminChannel:
			fmt.Printf("AdminAction: (%s)", adminEvent.Action)
		case gameEvent := <-e.GameChannel:
			fmt.Printf("GameAction: (%+v)\n", gameEvent.Action)
			e.PerformGameAction(gameEvent)
		}
	}
}
