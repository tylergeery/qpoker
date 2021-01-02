package gameplay

import (
	"database/sql"
	"fmt"
	"qpoker/cards"
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"qpoker/qutils"
	"qpoker/websocket/events"
	"time"
)

// GameState controls the game state returned to clients
type GameState struct {
	Manager *holdem.GameManager    `json:"manager"`
	Cards   map[int64][]cards.Card `json:"cards"`
}

// NewGameState returns the game state for clients
func NewGameState(manager *holdem.GameManager) GameState {
	return GameState{
		Manager: manager,
	}
}

// HoldemGameController is the GameController interface implementation for Holdem
type HoldemGameController struct {
	controller *Controller
	manager *holdem.GameManager
}

// Data returns Controller data object
func (c *HoldemGameController) Data() *Controller {
	return c.controller
}

// PerformGameAction advances to next hand
func (c *HoldemGameController) PerformGameAction(playerID int64, action interface{}, broadcast func(int64)) {
	complete, err := c.manager.PlayerAction(playerID, action.(holdem.Action))
	if err != nil {
		fmt.Printf("Error performing gameEvent: %+v, err: %s\n", action, err)
	}

	c.advance(complete, broadcast)
}

// Start holdem gameplay
func (c *HoldemGameController) Start(broadcast func(int64)) {
	c.nextHand()
	c.advance(false, broadcast)
}

func (c *HoldemGameController) advance(complete bool, broadcast func(int64)) {
	broadcast(c.controller.Game.ID)

	if complete {
		c.handComplete()
		broadcast(c.controller.Game.ID)
		return
	}

	if c.manager.IsAllIn() {
		complete := c.advanceAllIn()
		broadcast(c.controller.Game.ID)
		c.advance(complete, broadcast) // recursive call to finish game
	}
}

// UpdatePlayerChips updates players chips
func (c *HoldemGameController) UpdatePlayerChips(playerID, amount int64) {
	c.manager.AddChips(playerID, amount)
}

// AddPlayer adds new player
func (c *HoldemGameController) AddPlayer(player *models.Player) interface{} {
	gamePlayer := holdem.NewPlayer(player)
	c.manager.AddPlayer(gamePlayer)
	c.reloadPlayerStack(c.controller.Game, gamePlayer)

	return gamePlayer
}

// GetState returns visible game state for updating players
func (c *HoldemGameController) GetState(playerID int64) interface{} {
	state := NewGameState(c.manager)
	state.Cards = c.manager.GetVisibleCards(playerID)

	return state
}

// GetTimedOutGameEvent gets any needed overdue moves to auto advance
func (c *HoldemGameController) GetTimedOutGameEvent() (events.GameEvent, error) {
	lastMoveAt := c.manager.State.Table.ActiveAt
	currentTime := time.Now().Unix()
	allowedTime := qutils.ToInt(c.controller.Game.Options["decision_time"])
	if c.manager.Status != holdem.StatusActive {
		return events.GameEvent{}, fmt.Errorf("No idle event needed, game over")
	}

	if (currentTime - lastMoveAt) <= int64(allowedTime) {
		return events.GameEvent{}, fmt.Errorf("No idle event needed, still time")
	}

	action := holdem.Action{holdem.ActionFold, int64(0)}
	player := c.manager.State.Table.GetActivePlayer()

	if player.Options["can_check"] {
		action = holdem.Action{holdem.ActionCheck, int64(0)}
	}

	return events.GameEvent{
		GameID:   c.controller.Game.ID,
		PlayerID: player.ID,
		Action:   action,
	}, nil
}

func (c *HoldemGameController) handComplete() {
	waitSeconds := qutils.ToInt(c.controller.Game.Options["time_between_hands"])
	time.Sleep(time.Duration(waitSeconds) * time.Second)

	c.nextHand()
}

func (c *HoldemGameController) nextHand() {
	fmt.Printf("Advancing to next hand: %d\n", c.controller.Game.ID)
	err := c.manager.NextHand()
	if err != nil {
		fmt.Printf("Error advancing hand for admin: %s\n", err)
	}
}

func (c *HoldemGameController) advanceAllIn() bool {
	time.Sleep(time.Duration(2) * time.Second)
	fmt.Printf("Processing next round: %d\n", c.controller.Game.ID)

	complete, err := c.manager.ProcessAction()
	if err != nil {
		fmt.Printf("Error processing all in action: %s\n", err)
	}

	return complete
}

func (c *HoldemGameController) reloadPlayerStack(game *models.Game, player *holdem.Player) {
	// Search first for ending stack
	since := game.CreatedAt
	playerHand, err := models.GetGamePlayerHandForGameAndPlayer(game.ID, player.ID)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if playerHand.ID > int64(0) {
		player.Stack = playerHand.Starting
		if playerHand.Ending > -1 {
			player.Stack = playerHand.Ending
			since = playerHand.UpdatedAt
		}
	}

	// Add any chip requests since ending stack
	chipRequests, err := models.GetApprovedChipRequestsForGameAndPlayer(game.ID, player.ID, since, 1)
	if err != nil {
		return
	}

	if len(chipRequests) > 0 {
		player.Stack += chipRequests[0].Amount
	}
}
