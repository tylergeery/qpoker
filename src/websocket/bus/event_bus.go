package bus

import (
	"encoding/json"
	"fmt"
	"qpoker/models"
	"qpoker/websocket/connection"
	"qpoker/websocket/events"
	"qpoker/websocket/gameplay"
	"strconv"
	"time"
)

var eventBus *EventBus

// EventBus manages all server event action
type EventBus struct {
	games          map[int64]gameplay.GameController
	PlayerChannel  chan events.PlayerEvent
	GameChannel    chan events.GameEvent
	AdminChannel   chan events.AdminEvent
	MessageChannel chan events.MsgEvent
	VideoChannel   chan events.VideoEvent
}

// NewEventBus returns a new EventBus
func NewEventBus() *EventBus {
	return &EventBus{
		games:          map[int64]gameplay.GameController{},
		PlayerChannel:  make(chan events.PlayerEvent),
		GameChannel:    make(chan events.GameEvent),
		AdminChannel:   make(chan events.AdminEvent),
		MessageChannel: make(chan events.MsgEvent),
		VideoChannel:   make(chan events.VideoEvent),
	}
}

// StartEventBus creates and starts eventbus
func StartEventBus() *EventBus {
	eventBus = NewEventBus()

	go eventBus.ListenForEvents()
	go eventBus.IdleGameUpdates()

	return eventBus
}

func (e *EventBus) loadGameState(client *connection.Client) error {
	game, err := models.GetGameBy("id", client.GameID)
	if err != nil {
		return err
	}

	_, err = models.GetPlayerFromID(client.PlayerID)
	if err != nil {
		return err
	}

	controller, err := gameplay.GetGameController(game)
	if err != nil {
		return err
	}

	e.games[client.GameID] = controller

	return nil
}

func (e *EventBus) broadcast(gameID, playerID int64, broadcastEvent events.BroadcastEvent) {
	controller, ok := e.games[gameID]
	if !ok {
		return
	}

	state, err := json.Marshal(broadcastEvent)
	if err != nil {
		fmt.Printf("Error broadcasting game state: %s\n", err)
		return
	}

	clients := controller.Data().Clients
	for i := range clients {
		if clients[i].PlayerID != playerID {
			continue
		}

		clients[i].SendMessage(state)
	}
}

func (e *EventBus) handleAdminChipRequest(event events.AdminEvent) {
	controller := e.games[event.GameID]

	// Check if user has outstanding request
	requests := controller.Data().Requests
	for i := range requests {
		if requests[i].PlayerID == event.PlayerID {
			fmt.Printf("Only one chip request per player at a time: %d\n", event.PlayerID)
			return
		}
	}

	// Add request
	request := event.GetChipRequest()
	controller.Data().AddRequest(request)

	// Immediately approve for game owner
	if event.PlayerID == controller.Data().Game.OwnerID {
		event.Value = strconv.Itoa(int(request.PlayerID))
		e.handleAdminChipResponse(event)
		return
	}

	// Let owner know about request
	e.BroadcastRequests(controller)
}

func (e *EventBus) handleAdminChipResponse(event events.AdminEvent) {
	controller := e.games[event.GameID]
	id := event.Value.(string)
	approved := true

	// TODO: we can do better than have denied request be `-playerID`
	if id[0] == '-' {
		approved, id = false, id[1:]
	}

	playerID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Printf("Could not parse player ID: %s\n", id)
		return
	}

	// Remove request from pending
	chipRequest := controller.Data().RemovePlayerRequest(int64(playerID))
	if chipRequest == nil {
		fmt.Printf("Could not find chip request: %d %+v\n", playerID, controller.Data().Requests)
		return
	}

	// Approve chips and assign, if necessary
	chipRequest.Status = models.GameChipRequestStatusDenied
	if approved {
		chipRequest.Status = models.GameChipRequestStatusApproved
		controller.UpdatePlayerChips(chipRequest.PlayerID, chipRequest.Amount)
		e.BroadcastState(event.GameID)
	}

	// Save request
	err = chipRequest.Save()
	if err != nil {
		fmt.Printf("Error saving chip request (%s)\n", err)
	}

	// Let all players know about updated stack
	e.BroadcastRequests(controller)
}

func (e *EventBus) handleAdminEvent(event events.AdminEvent) {
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("Error: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	err := event.ValidateAuthorized(controller.Data().Game)
	if err != nil {
		return
	}

	switch event.Action {
	case connection.ClientAdminStart:
		controller.Start(e.BroadcastState)
		break
	case connection.ClientChipRequest:
		e.handleAdminChipRequest(event)
		break
	case connection.ClientChipResponse:
		e.handleAdminChipResponse(event)
		break
	default:
		fmt.Printf("Unknown admin event: %s\n", event.Action)
	}
}

func (e *EventBus) handleMessageEvent(event events.MsgEvent) {
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("Error: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	controller.Data().AddChat(event.GetChat())

	e.BroadcastMessages(controller)
}

func (e *EventBus) handleVideoEvent(event events.VideoEvent) {
	controller, ok := e.games[event.GameID]
	if !ok {
		fmt.Printf("VideoEvent error: Could not find controller for game id (%d)\n", event.GameID)
		return
	}

	broadcastEvent := events.NewBroadcastEvent(events.ActionVideo, event)
	clients := controller.Data().Clients

	for i := range clients {
		if clients[i].PlayerID == event.ToPlayerID {
			e.broadcast(clients[i].GameID, event.ToPlayerID, broadcastEvent)
		}
	}

	//e.BroadcastVideos(controller)
}

// BroadcastVideos broadcast client video states
func (e *EventBus) BroadcastVideos(controller gameplay.GameController) {
	clientVideos := map[int64]bool{}

	clients := controller.Data().Clients
	for i := range clients {
		clientVideos[clients[i].PlayerID] = true
	}

	broadcastEvent := events.NewBroadcastEvent(events.ActionVideo, clientVideos)
	for i := range clients {
		e.broadcast(clients[i].GameID, clients[i].PlayerID, broadcastEvent)
	}
}

// SetClient adds client to EventBus
func (e *EventBus) SetClient(client *connection.Client) {
	controller, ok := e.games[client.GameID]
	fmt.Println("Set client:", controller)
	if !ok {
		err := e.loadGameState(client)
		if err != nil {
			fmt.Printf("error loading game state: %s", err)
			return
		}

		controller = e.games[client.GameID]
	}

	player, err := models.GetPlayerFromID(client.PlayerID)
	if err != nil {
		return
	}

	controller.AddPlayer(player)
	controller.Data().AddClient(client)
	e.BroadcastState(client.GameID)

	if client.PlayerID == controller.Data().Game.OwnerID {
		e.BroadcastRequests(controller)
	}

	e.BroadcastMessages(controller)
	e.BroadcastVideos(controller)
}

// RemoveClient removes a client from EventBus
func (e *EventBus) RemoveClient(client *connection.Client) {
	controller, ok := e.games[client.GameID]
	if !ok {
		return
	}

	ok = controller.Data().RemoveClient(client)
	if !ok {
		return
	}

	e.BroadcastState(client.GameID)
	e.BroadcastVideos(controller)
}

// BroadcastRequests sends chip requests to game owner
func (e *EventBus) BroadcastRequests(controller gameplay.GameController) {
	game := controller.Data().Game
	broadcastEvent := events.NewBroadcastEvent(events.ActionAdmin, map[string][]*models.GameChipRequest{
		"requests": controller.Data().Requests,
	})

	e.broadcast(game.ID, game.OwnerID, broadcastEvent)
}

// BroadcastMessages sends chip requests to game owner
func (e *EventBus) BroadcastMessages(controller gameplay.GameController) {
	clients := controller.Data().Clients
	event := events.NewBroadcastEvent(events.ActionMsg, controller.Data().Chats)

	for i := range clients {
		e.broadcast(clients[i].GameID, clients[i].PlayerID, event)
	}
}

// BroadcastState sends game state to all clients
func (e *EventBus) BroadcastState(gameID int64) {
	controller, ok := e.games[gameID]
	if !ok {
		return
	}

	clients := controller.Data().Clients
	for i := range clients {
		state := controller.GetState(clients[i].PlayerID)
		broadcastEvent := events.NewBroadcastEvent(events.ActionGame, state)
		gameState, _ := json.Marshal(broadcastEvent)
		clients[i].SendMessage(gameState)
	}
}

// PerformGameAction calls perform game action on controller and broadcasts updates
func (e *EventBus) PerformGameAction(gameEvent events.GameEvent) {
	controller, ok := e.games[gameEvent.GameID]
	if !ok {
		return
	}

	controller.PerformGameAction(gameEvent.PlayerID, gameEvent.Action, e.BroadcastState)
}

// ListenForEvents starts the event bus waiting for channel events
func (e *EventBus) ListenForEvents() {
	for {
		select {
		case playerEvent := <-e.PlayerChannel:
			fmt.Printf("PlayerAction: (%+v)\n", playerEvent)
			playerEventMap := map[string]func(*connection.Client){
				events.ActionPlayerRegister: e.SetClient,
				events.ActionPlayerLeave:    e.RemoveClient,
			}
			playerEventMap[playerEvent.Action](playerEvent.Client)
		case adminEvent := <-e.AdminChannel:
			fmt.Printf("AdminAction: (%+v)\n", adminEvent)
			e.handleAdminEvent(adminEvent)
		case gameAction := <-e.GameChannel:
			fmt.Printf("GameAction: (%+v)\n", gameAction)
			e.PerformGameAction(gameAction)
		case msgAction := <-e.MessageChannel:
			fmt.Printf("MsgAction: (%+v)\n", msgAction)
			e.handleMessageEvent(msgAction)
		case videoAction := <-e.VideoChannel:
			fmt.Printf("VideoAction: (%+v)\n", videoAction)
			e.handleVideoEvent(videoAction)
		}
	}
}

// IdleGameUpdates handles idle time updates for games
func (e *EventBus) IdleGameUpdates() {
	for {
		for _, controller := range e.games {
			gameEvent, err := controller.GetTimedOutGameEvent()
			if err == nil {
				e.GameChannel <- gameEvent
			}
		}

		// TODO: heuristic sleep?
		time.Sleep(200 * time.Millisecond)
	}
}
