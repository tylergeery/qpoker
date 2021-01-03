package gameplay

import (
	"fmt"
	"qpoker/cards"
	"qpoker/cards/games/hearts"
	"qpoker/models"
	"qpoker/websocket/events"
)

// HeartsGameState controls the game state returned to clients
type HeartsGameState struct {
	Manager *hearts.GameManager    `json:"manager"`
	Cards   map[int64][]cards.Card `json:"cards"`
}

// NewHeartsGameState returns the game state for clients
func NewHeartsGameState(manager *hearts.GameManager) HeartsGameState {
	return HeartsGameState{
		Manager: manager,
	}
}

// HeartsGameController is the hearts implementation of GameController interface
type HeartsGameController struct {
	controller *Controller
	manager    *hearts.GameManager
}

// Data returns Controller data object
func (c *HeartsGameController) Data() *Controller {
	return c.controller
}

// PerformGameAction advances to next hand
func (c *HeartsGameController) PerformGameAction(playerID int64, action interface{}, broadcast func(int64)) {
	complete, err := c.manager.PlayerAction(playerID, action.(hearts.Action))
	if err != nil {
		fmt.Printf("Error performing gameEvent: %+v, err: %s\n", action, err)
	}

	c.advance(complete, broadcast)
}

// Start hearts gameplay
func (c *HeartsGameController) Start(broadcast func(int64)) {
	c.manager.NextHand()
	c.advance(false, broadcast)
}

// UpdatePlayerChips updates players chips
func (c *HeartsGameController) UpdatePlayerChips(playerID, amount int64) {
	// unneccessary for hearts
}

// AddPlayer adds new player
func (c *HeartsGameController) AddPlayer(player *models.Player) interface{} {
	gamePlayer := hearts.NewPlayer(player)
	c.manager.AddPlayer(gamePlayer)

	return gamePlayer
}

// GetState returns visible game state for updating players
func (c *HeartsGameController) GetState(playerID int64) interface{} {
	state := NewHeartsGameState(c.manager)
	state.Cards = c.manager.GetVisibleCards(playerID)

	return state
}

// GetTimedOutGameEvent gets a moved thats over time limit
func (c *HeartsGameController) GetTimedOutGameEvent() (events.GameEvent, error) {
	return events.GameEvent{}, fmt.Errorf("No idle game event needed")
}

func (c *HeartsGameController) advance(complete bool, broadcast func(int64)) {
	broadcast(c.controller.Game.ID)
}
