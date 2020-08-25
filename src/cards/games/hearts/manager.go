package hearts

import (
	"fmt"
	"qpoker/cards"
	"qpoker/models"
)

const (
	// MaxPlayerCount is the max amount of players for hearts
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

// GameOptions holds game options
type GameOptions struct {
	Capacity     int   `json:"capacity"`
	BigBlind     int64 `json:"big_blind"`
	DecisionTime int   `json:"decision_time"`
}

// NewGameManager returns a new GameManager
func NewGameManager(gameID int64, players []*Player, options GameOptions) (*GameManager, error) {
	if len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	table := NewTable(4, players)
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

	if g.State.skipPass() {
		g.UpdateStatus(StatusActive)
	}

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
	playersRemaining := g.State.Table.GetAllPlayers()

	for i := range playersRemaining {
		if len(playersRemaining[i].Cards) > 0 {
			return false
		}
	}

	return true
}

// PlayerPass passes player cards
func (g *GameManager) PlayerPass(playerID int64, cards []cards.Card) error {
	if len(cards) != 3 {
		return fmt.Errorf("Invalid pass, requires 3 cards: (%d) (%+v)", playerID, cards)
	}

	g.State.addPass(playerID, cards)

	if g.State.passesComplete() {
		g.UpdateStatus(StatusActive)
	}

	return nil
}

func (g *GameManager) playerPlay(action Action) error {
	player := g.State.Table.GetActivePlayer()

	if len(player.Cards) == 0 {
		return fmt.Errorf("Player (%d) cannot play, has no cards", player.ID)
	}

	if len(g.State.Board) == 4 {
		err := g.CleanPile()
		if err != nil {
			return err
		}
	}

	found := false
	for i := range player.Cards {
		if player.Cards[i].ToString() == action.Card.ToString() {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Card (%s) does not exist for player (%d)", action.Card.ToString(), player.ID)
	}

	g.State.playerPlay(player.ID, action.Card)
	return nil
}

// PlayerAction performs an action for player
func (g *GameManager) PlayerAction(playerID int64, action Action) (bool, error) {
	if g.isComplete() {
		return false, fmt.Errorf("Game is already complete")
	}

	if g.Status != StatusActive {
		return false, fmt.Errorf("Game is not yet active: %s", g.Status)
	}

	if playerID != g.State.Table.GetActivePlayer().ID {
		return false, fmt.Errorf("User (%d) must wait for player (%d) to act", playerID, g.State.Table.GetActivePlayer().ID)
	}

	actionMap := map[string]func(action Action) error{
		ActionPlay: g.playerPlay,
	}

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

	if g.Status == StatusActive {
		g.State.Table.ActivateNextPlayer(g.GetPlayerActions)
	}

	return false, nil
}

// GetPlayerActions returns allowed active player actions
func (g *GameManager) GetPlayerActions() map[string]bool {
	return map[string]bool{
		"can_play": true,
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
	for _, player := range g.State.Table.GetAllPlayers() {
		g.gamePlayerHands[player.ID] = &models.GamePlayerHand{
			GameHandID: g.gameHand.ID,
			Cards:      g.cardsToStringArray(player.Cards),
			PlayerID:   player.ID,
			Starting:   player.Score,
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
	g.State.CleanPile(false)

	players := g.State.Table.GetAllPlayers()
	g.gameHand.Payouts = g.State.PointTotals(players)

	err := g.gameHand.Save()
	if err != nil {
		return err
	}

	for _, player := range players {
		hand, ok := g.gamePlayerHands[player.ID]
		if !ok {
			continue
		}

		if amount, ok := g.gameHand.Payouts[player.ID]; ok {
			player.Score += amount
		}

		hand.Ending = player.Score
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
		if len(g.State.Table.GetAllPlayers()) > 3 {
			g.Status = status
		}
		break
	case g.Status == StatusReady && status == StatusInit:
		if len(g.State.Table.GetAllPlayers()) <= 3 {
			g.Status = status
		}
	case g.Status == StatusReady && status == StatusActive:
		if len(g.State.Table.GetAllPlayers()) > 3 {
			g.Status = status
		}
		break
	}
}

// CleanPile collects cards and adds them to player's pile
func (g *GameManager) CleanPile() error {
	return g.State.CleanPile(true)
}
