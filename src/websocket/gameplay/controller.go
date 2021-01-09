package gameplay

import (
	"fmt"
	"qpoker/cards/games/hearts"
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"qpoker/websocket/connection"
	"qpoker/websocket/events"
)

// Controller object holds controller data
type Controller struct {
	Clients  []*connection.Client
	Game     *models.Game
	Requests []*models.GameChipRequest
	Chats    []events.Chat
}

// AddRequest adds chip request
func (c *Controller) AddRequest(request *models.GameChipRequest) {
	c.Requests = append(c.Requests, request)
}

// RemovePlayerRequest removes chip request for player
func (c *Controller) RemovePlayerRequest(playerID int64) *models.GameChipRequest {
	for i := range c.Requests {
		if c.Requests[i].PlayerID == playerID {
			chipRequest := c.Requests[i]
			c.Requests = append(c.Requests[:i], c.Requests[i+1:]...)
			return chipRequest
		}
	}

	return nil
}

// AddChat adds chat to history
func (c *Controller) AddChat(chat events.Chat) {
	c.Chats = append(c.Chats, chat)
	if l := len(c.Chats); l > 100 {
		c.Chats = c.Chats[l-100:]
	}
}

// AddClient records new client
func (c *Controller) AddClient(client *connection.Client) {
	c.Clients = append(c.Clients, client)
}

// RemoveClient removes client if exists
func (c *Controller) RemoveClient(client *connection.Client) bool {
	for i := range c.Clients {
		if c.Clients[i] == client {
			c.Clients = append(c.Clients[:i], c.Clients[i+1:]...)
			return true
		}
	}

	return false
}

// GameController handles logic for sending/receiving game events
// and controlling game advancement calls to GameManager
type GameController interface {
	Data() *Controller
	PerformGameAction(playerID int64, action interface{}, broadcast func(int64))
	Start(broadcast func(int64))
	Pause(resume bool)
	UpdatePlayerChips(playerID, amount int64)
	AddPlayer(player *models.Player) interface{}
	GetState(playerID int64) interface{}
	GetTimedOutGameEvent() (events.GameEvent, error)
}

func buildHoldemController(game *models.Game) (GameController, error) {
	players := []*holdem.Player{}
	manager, err := holdem.NewGameManager(game.ID, players, holdem.NewGameOptions(game.Options))
	if err != nil {
		return nil, err
	}
	controller := &Controller{
		[]*connection.Client{}, game,
		[]*models.GameChipRequest{}, []events.Chat{},
	}

	return &HoldemGameController{controller, manager}, nil
}

func buildHeartsController(game *models.Game) (GameController, error) {
	players := []*hearts.Player{}
	manager, err := hearts.NewGameManager(game.ID, players, hearts.NewGameOptions(game.Options))
	if err != nil {
		return nil, err
	}
	controller := &Controller{
		[]*connection.Client{}, game,
		[]*models.GameChipRequest{}, []events.Chat{},
	}

	return &HeartsGameController{controller, manager}, nil
}

// GetGameController returns the appropriate interface implementation of the GameController
func GetGameController(game *models.Game) (GameController, error) {

	switch game.Status {
	case models.GameStatusComplete:
		return nil, fmt.Errorf("Game is already complete: %d", game.ID)
	}

	// TODO: do better than ids
	switch game.GameTypeID {
	case int64(1):
		return buildHoldemController(game)
	case int64(2):
		return buildHeartsController(game)
	default:
		return nil, fmt.Errorf("Unknown gameTypeId: %d", game.GameTypeID)
	}
}
