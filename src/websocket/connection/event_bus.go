package connection

import (
	"encoding/json"
	"fmt"
	"qpoker/cards"
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"time"
)

var eventBus *EventBus

// GameState controls the game state returned to clients
type GameState struct {
	Manager *holdem.GameManager    `json:"manager"`
	Cards   map[int64][]cards.Card `json:"cards"`
}

// ChipRequest holds a chips request
type ChipRequest struct {
	ID       string `json:"id"`
	PlayerID int64  `json:"player_id"`
	Amount   int64  `json:"amount"`
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
	clients  []*Client
	manager  *holdem.GameManager
	game     *models.Game
	requests []ChipRequest
	chats    []Chat
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

	player, err := models.GetPlayerFromID(client.PlayerID)
	if err != nil {
		return err
	}

	// TODO: Check if game is complete
	// TODO: pull latest game hand, recreate state from hand
	gamePlayer := holdem.NewPlayer(player)
	players := []*holdem.Player{gamePlayer}
	manager, err := holdem.NewGameManager(game.ID, players, game.Options)
	if err != nil {
		return err
	}

	controller := &GameController{[]*Client{}, manager, game, []ChipRequest{}}
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

func (e *EventBus) handleAdminChipRequest(event AdminEvent) {
	controller := e.games[event.GameID]
	request := event.GetChipRequest()
	controller.requests = append(controller.requests, request)

	if event.PlayerID == controller.game.OwnerID {
		event.Value = request.ID
		e.handleAdminChipResponse(event)
		return
	}

	e.BroadcastRequests(controller)
}

func (e *EventBus) handleAdminChipResponse(event AdminEvent) {
	controller := e.games[event.GameID]
	id := event.Value.(string)
	approved := true
	if id[0] == '-' {
		approved, id = false, id[1:]
	}

	found, chipRequest := false, ChipRequest{}
	for i := range controller.requests {
		if controller.requests[i].ID == id {
			found, chipRequest = true, controller.requests[i]
			controller.requests = append(controller.requests[:i], controller.requests[i+1:]...)
			break
		}
	}

	if !found {
		fmt.Printf("Could not find chip request: %s %+v\n", id, controller.requests)
		return
	}

	if approved {
		controller.manager.AddChips(chipRequest.PlayerID, chipRequest.Amount)
		e.BroadcastState(event.GameID)
	}

	e.BroadcastRequests(controller)
}

func (e *EventBus) advanceNextHand(event AdminEvent) {
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("Error advancing hand: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	err := controller.manager.NextHand()
	if err != nil {
		fmt.Printf("Error advancing hand for admin: %s\n", err)
		return
	}

	e.BroadcastState(event.GameID)
}

func (e *EventBus) handleAdminEvent(event AdminEvent) {
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("Error: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	err := event.ValidateAuthorized(controller.game)
	if err != nil {
		return
	}

	switch event.Action {
	case ClientAdminStart:
		e.advanceNextHand(event)
		break
	case ClientChipRequest:
		e.handleAdminChipRequest(event)
		break
	case ClientChipResponse:
		e.handleAdminChipResponse(event)
		break
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

	player, err := models.GetPlayerFromID(client.PlayerID)
	if err != nil {
		return
	}

	gamePlayer := holdem.NewPlayer(player)
	_ = controller.manager.AddPlayer(gamePlayer)
	controller.clients = append(controller.clients, client)

	e.BroadcastState(client.GameID)

	if client.PlayerID == controller.game.OwnerID {
		e.BroadcastRequests(controller)
	}
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
			return
		}
	}

	fmt.Println("Broadcasting client state: remove client")
	e.BroadcastState(client.GameID)
}

// BroadcastRequests sends chip requests to game owner
func (e *EventBus) BroadcastRequests(controller *GameController) {
	// TODO: include relevant cards
	broadcastEvent := NewBroadcastEvent(ActionAdmin, map[string][]ChipRequest{
		"requests": controller.requests,
	})

	fmt.Printf("broadcast %+v %+v\n", controller.game, broadcastEvent)
	e.broadcast(controller.game.ID, controller.game.OwnerID, broadcastEvent)
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
		time.AfterFunc(
			controller.game.Options.TimeBetweenHands*time.Second,
			func() {
				e.advanceNextHand(AdminEvent{GameID: gameEvent.GameID})
			},
		)
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
			fmt.Printf("AdminAction: (%+v)\n", adminEvent)
			e.handleAdminEvent(adminEvent)
		case gameAction := <-e.GameChannel:
			fmt.Printf("GameAction: (%+v)\n", gameAction)
			e.PerformGameAction(gameAction)
		}
	}
}
