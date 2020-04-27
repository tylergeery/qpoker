package connection

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"qpoker/cards"
	"qpoker/cards/games/holdem"
	"qpoker/models"
)

var eventBus *EventBus

// GameState controls the game state returned to clients
type GameState struct {
	Manager *holdem.GameManager     `json:"manager"`
	Cards   map[int64][]cards.Card  `json:"cards"`
	Players map[int64]holdem.Player `json:"players"`
}

// ChipRequest holds a chips request
type ChipRequest struct {
	ID string `json:"id"`
	PlayerID int64 `json:"player_id"`
	Amount int64 `json:"amount"`
}

// NewGameState returns the game state for clients
func NewGameState(manager *holdem.GameManager) GameState {
	return GameState{
		Manager: manager,
		Cards:   manager.GetVisibleCards(),
	}
}

// GameController handles logic for sending/receiving game events
type GameController struct {
	clients []*Client
	manager *holdem.GameManager
	game    *models.Game
	requests []
}

// EventBus manages all server event action
type EventBus struct {
	games          map[int64]*GameController
	PlayerChannel  chan PlayerEvent
	GameChannel    chan GameEvent
	AdminChannel   chan AdminEvent
	MessageChannel chan MsgEvent
}

// NewEventBus returns a new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		games:          map[int64]*GameController{},
		PlayerChannel:  make(chan PlayerEvent),
		GameChannel:    make(chan GameEvent),
		AdminChannel:   make(chan AdminEvent),
		MessageChannel: make(chan MsgEvent),
	}
}

// StartEventBus creates and starts eventbus
func StartEventBus() *EventBus {
	eventBus = NewEventBus()

	go eventBus.ListenForEvents()

	return eventBus
}

func (e *EventBus) reloadGameState(client *Client) error {
	game, err := models.GetGameBy("id", client.GameID)
	if err != nil {
		return err
	}

	// TODO: Check if game is complete
	// TODO: pull latest game hand, recreate state from hand
	players := []*holdem.Player{&holdem.Player{ID: client.PlayerID}}
	manager, err := holdem.NewGameManager(game.ID, players, game.Options)
	if err != nil {
		return err
	}

	controller := &GameController{[]*Client{}, manager, game}
	e.games[client.GameID] = controller

	return nil
}

func (e *EventBus) broadcast(gameID, playerID int64, broadcastEvent BroadcastEvent) {
	controller, ok := e.games[gameID]
	if !ok {
		return
	}

	state, err := json.Marshal(broadcastEvent)
	if err != nil {
		fmt.Printf("Error broadcasting game state: %s\n", err)
		return
	}

	for i := range controller.clients {
		if controller.clients[i].PlayerID != playerID {
			continue
		}

		controller.clients[i].SendMessage(state)
	}
}

func (e *EventBus) handleAdminEvent(event AdminEvent) {
	fmt.Println("handleAdminEvent")
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("Error: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	fmt.Printf("validateAuthorized: %d, %+v\n", event.PlayerID, controller.game)
	err := event.ValidateAuthorized(controller.game)
	if err != nil {
		return
	}

	fmt.Println("event.Action", event.Action)
	switch event.Action {
	case ClientAdminStart:
		err = controller.manager.NextHand()
		if err != nil {
			fmt.Printf("Error starting game for admin: %s\n", err)
			return
		}

		e.BroadcastState(event.GameID)
	case ClientChipRequest:
		broadcastEvent := NewBroadcastEvent(
			ActionAdmin,
			map[string]interface{}{
				"player": controller.manager.GetPlayer(event.PlayerID),
				"amount": int64(event.Value.(float64)),
			},
		)

		fmt.Printf("broadcast %+v %+v\n", controller.game, broadcastEvent)
		e.broadcast(controller.game.ID, controller.game.OwnerID, broadcastEvent)
	case ClientChipResponse:
		amount := int64(event.Value.(float64))
		controller.manager.GetPlayer(event.PlayerID).Stack += amount

		e.BroadcastState(event.GameID)
	default:
		fmt.Printf("Unknown admin event: %s\n", event.Action)
	}
}

// SetClient adds client to EventBus
func (e *EventBus) SetClient(client *Client) {
	controller, ok := e.games[client.GameID]
	fmt.Println("Set client:", controller)
	if !ok {
		err := e.reloadGameState(client)
		if err != nil {
			fmt.Printf("error reloading game state: %s", err)
			return
		}

		controller = e.games[client.GameID]
	}

	client.GameChannel = e.GameChannel
	client.AdminChannel = e.AdminChannel
	client.MessageChannel = e.MessageChannel

	_ = controller.manager.State.Table.AddPlayer(&holdem.Player{ID: client.PlayerID})
	controller.clients = append(controller.clients, client)

	fmt.Println("Broadcasting client state: add client")
	e.BroadcastState(client.GameID)
}

// RemoveClient removes a client from EventBus
func (e *EventBus) RemoveClient(client *Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		return
	}

	// TODO: this is broken
	for i := range controller.clients {
		if controller.clients[i] == client {
			controller.clients = append(controller.clients[:i], controller.clients[i+1:]...)
			return
		}
	}

	fmt.Println("Broadcasting client state: remove client")
	e.BroadcastState(client.GameID)
}

// BroadcastState sends game state to all clients
func (e *EventBus) BroadcastState(gameID int64) {
	controller, ok := e.games[gameID]
	if !ok {
		return
	}

	// TODO: include relevant cards
	broadcastEvent := NewBroadcastEvent(ActionGame, NewGameState(controller.manager))
	state, err := json.Marshal(broadcastEvent)
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
			fmt.Printf("AdminAction: (%s)\n", adminEvent.Action)
			e.handleAdminEvent(adminEvent)
		case gameAction := <-e.GameChannel:
			fmt.Printf("GameAction: (%+v)\n", gameAction.Action)
			e.PerformGameAction(gameAction)
		}
	}
}
