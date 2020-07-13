package holdem

import (
	"fmt"
	"qpoker/cards"
	"qpoker/models"
)

const (
	// MaxPlayerCount is the max amount of players for HoldEm
	MaxPlayerCount = 4

	// StatusInit init
	StatusInit = "init"
	// StatusReady ready
	StatusReady = "ready"
	// StatusActive active
	StatusActive = "active"
)

// GameManager holds game state and manages game flow
type GameManager struct {
	GameID int64   `json:"game_id"`
	State  *Hearts `json:"state"`
	Status string  `json:"status"`

	gameHand        *models.GameHand
	gamePlayerHands map[int64]*models.GamePlayerHand
}

// NewGameManager returns a new GameManager
func NewGameManager(gameID int64, players []*Player, options models.GameOptions) (*GameManager, error) {
	if len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	table := NewTable(players)
	gm := &GameManager{
		GameID: gameID,
		State:  NewHearts(table, ""), // TODO
		Status: StatusInit,
	}

	return gm, nil
}

// AddPlayer adds player to game
func (g *GameManager) AddPlayer(player *Player) error {
	err := g.State.Table.AddPlayer(player)

	g.UpdateStatus(StatusReady)

	return err
}

// RemovePlayer removes player from game
func (g *GameManager) RemovePlayer(playerID int64) error {
	err := g.State.Table.RemovePlayer(playerID)

	g.UpdateStatus(StatusInit)

	return err
}

// NextHand moves game manager on to next hand
func (g *GameManager) NextHand() error {
	g.Status = StatusReady

	err := g.State.Deal()
	if err != nil {
		return err
	}

	// Save hand and player hands
	err = g.StartHand()
	if err != nil {
		return err
	}

	g.UpdateStatus(StatusActive)
	g.ProcessAction()

	return nil
}

func (g *GameManager) cardsToStringArray(cardObjects []cards.Card) []string {
	stringCards := make([]string, len(cardObjects))

	for i, c := range cardObjects {
		stringCards[i] = c.ToString()
	}

	return stringCards
}

func (g *GameManager) isComplete() bool {
	playersRemaining := g.State.Table.GetPlayers()

	for i := range playersRemaining {
		if len(playersRemaining[i].Cards) > 0 {
			return false
		}
	}

	return true
}

// PlayerAction performs an action for player
func (g *GameManager) PlayerAction(playerID int64, action Action) (bool, error) {
	if g.isComplete() {
		return false, fmt.Errorf("Game is already complete")
	}

	if playerID != g.State.Table.GetActivePlayer().ID {
		return false, fmt.Errorf("User (%d) must wait for player (%d) to act", playerID, g.State.Table.GetActivePlayer().ID)
	}

	actionMap := map[string]func(action Action) error{}

	err := actionMap[action.Name](action)
	if err != nil {
		return false, err
	}

	return g.ProcessAction()
}

// ProcessAction handles the post-processing of an action
func (g *GameManager) ProcessAction() (bool, error) {
	if g.isComplete() {
		err := g.EndHand()
		g.State.Table.GetActivePlayer().SetPlayerActions(nil)
		return true, err
	}

	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	return false, nil
}

// GetPlayerActions returns allowed active player actions
func (g *GameManager) GetPlayerActions() map[string]bool {
	return map[string]bool{
		// "can_play": g.canPlay(),
	}
}

// StartHand saves the initial game state to the DB
func (g *GameManager) StartHand() error {
	g.gameHand = &models.GameHand{GameID: g.GameID}
	err := g.gameHand.Save()
	if err != nil {
		return err
	}

	g.gamePlayerHands = map[int64]*models.GamePlayerHand{}
	for _, player := range g.State.Table.GetPlayers() {
		g.gamePlayerHands[player.ID] = &models.GamePlayerHand{
			GameHandID:    g.gameHand.ID,
			Cards:         g.cardsToStringArray(player.Cards),
			PlayerID:      player.ID,
			StartingStack: int64(player.Score),
		}
		err = g.gamePlayerHands[player.ID].Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// EndHand saves the ending game state to the DB
func (g *GameManager) EndHand() error {
	// g.gameHand.Payouts = payouts

	// err := g.gameHand.Save()
	// if err != nil {
	// 	return err
	// }

	for _, player := range g.State.Table.Players {
		hand, ok := g.gamePlayerHands[player.ID]
		if !ok {
			continue
		}

		// if amount, ok := payouts[player.ID]; ok {
		// 	player.Stack += amount
		// }

		err := hand.Save()
		if err != nil {
			fmt.Printf("Error saving user hand: %s\n", err)
		}
	}

	return nil
}

// UpdateStatus updates game status to most appropriate
func (g *GameManager) UpdateStatus(status string) {
	switch {
	case g.Status == StatusInit && status == StatusReady:
		if len(g.State.Table.GetPlayers()) > 1 {
			g.Status = status
		}
		break
	case g.Status == StatusReady && status == StatusInit:
		if len(g.State.Table.GetActivePlayers()) <= 1 {
			g.Status = status
		}
	case g.Status == StatusReady && status == StatusActive:
		if len(g.State.Table.GetActivePlayers()) > 1 {
			g.Status = status
		}
		break
	default:
		break
	}
}

// GetPlayer returns a player from a table
func (g *GameManager) GetPlayer(playerID int64) *Player {
	for i := range g.State.Table.Players {
		if g.State.Table.Players[i] == nil {
			continue
		}

		if playerID == g.State.Table.Players[i].ID {
			return g.State.Table.Players[i]
		}
	}

	return nil
}
