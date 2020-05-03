package holdem

import (
	"fmt"
	"qpoker/cards"
	"qpoker/models"
	"qpoker/utils"
)

const (
	// MaxPlayerCount is the max amount of players for HoldEm
	MaxPlayerCount = 12

	// StatusInit init
	StatusInit = "init"
	// StatusReady ready
	StatusReady = "ready"
	// StatusActive active
	StatusActive = "active"
)

// GameManager holds game state and manages game flow
type GameManager struct {
	GameID   int64   `json:"game_id"`
	State    *HoldEm `json:"state"`
	Pot      *Pot    `json:"pot"`
	BigBlind int64   `json:"big_blind"`
	Status   string  `json:"status"`

	gameHand        *models.GameHand
	gamePlayerHands map[int64]*models.GamePlayerHand
}

// NewGameManager returns a new GameManager
func NewGameManager(gameID int64, players []*Player, options models.GameOptions) (*GameManager, error) {
	if len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	if options.Capacity > 0 && len(players) > options.Capacity {
		return nil, fmt.Errorf("Player count (%d) is greater than capacity (%d)", len(players), options.Capacity)
	}

	table := NewTable(options.Capacity, players)
	gm := &GameManager{
		GameID:   gameID,
		State:    NewHoldEm(table),
		BigBlind: options.BigBlind,
		Status:   StatusInit,
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

// AddChips adds chips to Player
func (g *GameManager) AddChips(playerID, amount int64) {
	g.GetPlayer(playerID).Stack += amount
	g.UpdateStatus(StatusReady)
}

// NextHand moves game manager on to next hand
func (g *GameManager) NextHand() error {
	g.UpdateStatus(StatusActive)

	err := g.State.Deal()
	if err != nil {
		return err
	}

	// TODO: pull latest game options

	// Save hand and player hands
	err = g.StartHand()
	if err != nil {
		return err
	}

	g.Pot = NewPot(g.State.Table.GetActivePlayers())

	// little
	g.State.Table.GetActivePlayer().LittleBlind = true
	g.playerBet(NewActionBet(utils.MinInt64(g.BigBlind/2, g.State.Table.GetActivePlayer().Stack)))
	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	// big blind
	g.State.Table.GetActivePlayer().BigBlind = true
	g.playerBet(NewActionBet(utils.MinInt64(g.BigBlind, g.State.Table.GetActivePlayer().Stack)))
	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

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
	playersRemaining := g.State.Table.GetActivePlayers()

	if len(playersRemaining) <= 1 {
		return true
	}

	if g.isRoundComplete() && g.State.State == StateRiver {
		return true
	}

	return false
}

func (g *GameManager) isRoundComplete() bool {
	playersRemaining := g.State.Table.GetActivePlayers()
	nextPlayer := playersRemaining[0]

	// Check if everyone has called/checked the next player
	for _, player := range playersRemaining {
		// Everyone gets a chance to play
		if player.State == PlayerStateInit {
			return false
		}

		// If anybody besides the next player bet, keep going
		if player.ID != nextPlayer.ID && player.State == PlayerStateBet {
			return false
		}
	}

	if g.State.State != StateDeal {
		return true
	}

	if !nextPlayer.BigBlind {
		return true
	}

	return g.Pot.MaxBet() != g.BigBlind
}

func (g *GameManager) canBet() bool {
	if g.State.Table.GetActivePlayer().Stack <= int64(0) {
		return false
	}

	return true
}

func (g *GameManager) canCall() bool {
	if !g.canBet() {
		return false
	}

	return g.Pot.MaxBet() > g.Pot.PlayerBets[g.State.Table.GetActivePlayer().ID]
}

func (g *GameManager) canCheck() bool {
	return g.Pot.MaxBet() == g.Pot.PlayerBets[g.State.Table.GetActivePlayer().ID]
}

func (g *GameManager) canFold() bool {
	return !g.canCheck()
}

func (g *GameManager) playerBet(action Action) error {
	if !g.canBet() {
		return fmt.Errorf("Cannot bet")
	}

	// TODO: validate bet amount, don't forget about little blind
	player := g.State.Table.GetActivePlayer()
	amount := utils.MinInt64(action.Amount, player.Stack)

	g.Pot.AddBet(player.ID, amount)
	player.Stack -= amount
	player.State = PlayerStateBet

	return nil
}

func (g *GameManager) playerCall(action Action) error {
	if !g.canCall() {
		return fmt.Errorf("Cannot call")
	}

	player := g.State.Table.GetActivePlayer()
	amount := g.Pot.MaxBet() - g.Pot.PlayerBets[player.ID]

	g.Pot.AddBet(player.ID, amount)
	player.Stack -= amount
	player.State = PlayerStateCall

	return nil
}

func (g *GameManager) playerCheck(action Action) error {
	if !g.canCheck() {
		return fmt.Errorf("Cannot check")
	}

	g.State.Table.GetActivePlayer().State = PlayerStateCheck

	return nil
}

func (g *GameManager) playerFold(action Action) error {
	if !g.canFold() {
		return fmt.Errorf("Cannot fold")
	}

	fmt.Printf("setting player state: %s\n", PlayerStateFold)
	g.State.Table.GetActivePlayer().State = PlayerStateFold

	return nil
}

// PlayerAction performs an action for player
func (g *GameManager) PlayerAction(playerID int64, action Action) (complete bool, err error) {
	if g.isComplete() {
		err = fmt.Errorf("Game is already complete")
		return
	}

	if playerID != g.State.Table.GetActivePlayer().ID {
		err = fmt.Errorf("User (%d) must wait for player (%d) to act", playerID, g.State.Table.GetActivePlayer().ID)
		return
	}

	actionMap := map[string]func(action Action) error{
		ActionBet:   g.playerBet,
		ActionCall:  g.playerCall,
		ActionCheck: g.playerCheck,
		ActionFold:  g.playerFold,
	}

	err = actionMap[action.Name](action)
	if err != nil {
		return
	}

	complete = g.isComplete()
	if complete {
		err = g.EndHand()
		g.State.Table.GetActivePlayer().SetPlayerActions(nil)
		return
	}

	if g.isRoundComplete() {
		g.State.Advance()
		g.State.Table.NextRound(g.GetPlayerActions)
		g.Pot.ClearBets()
		return
	}

	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	return
}

// GetPlayerActions returns allowed active player actions
func (g *GameManager) GetPlayerActions() map[string]bool {
	return map[string]bool{
		"can_bet":   g.canBet(),
		"can_call":  g.canCall(),
		"can_check": g.canCheck(),
		"can_fold":  g.canFold(),
	}
}

// ShowVisibleCards marks remaining players as visible cards
func (g *GameManager) ShowVisibleCards() {
	activePlayers := g.State.Table.GetActivePlayers()

	if len(activePlayers) <= 1 {
		return
	}

	for _, player := range activePlayers {
		player.CardsVisible = true
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
	for _, player := range g.State.Table.GetActivePlayers() {
		g.gamePlayerHands[player.ID] = &models.GamePlayerHand{
			GameHandID:    g.gameHand.ID,
			Cards:         g.cardsToStringArray(player.Cards),
			PlayerID:      player.ID,
			StartingStack: player.Stack,
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
	g.ShowVisibleCards()
	payouts := g.Pot.GetPayouts(g.State.GetWinningIDs())

	g.gameHand.Board = g.cardsToStringArray(g.State.Board)
	g.gameHand.Payouts = payouts
	g.gameHand.Bets = g.Pot.PlayerTotals

	err := g.gameHand.Save()
	if err != nil {
		return err
	}

	for _, player := range g.State.Table.Players {
		if player == nil {
			continue
		}

		hand, ok := g.gamePlayerHands[player.ID]
		if !ok {
			continue
		}

		hand.EndingStack = player.Stack
		if amount, ok := payouts[player.ID]; ok {
			hand.EndingStack += amount
		}

		err = hand.Save()
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
		if len(g.State.Table.GetAllPlayers()) > 1 {
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

// GetVisibleCards returns client representation of cards for those visible
func (g *GameManager) GetVisibleCards(playerID int64) map[int64][]cards.Card {
	visibleCards := map[int64][]cards.Card{}

	for _, player := range g.State.Table.Players {
		if player == nil {
			continue
		}

		if player.Cards == nil {
			continue
		}

		if player.ID == playerID || player.CardsVisible {
			visibleCards[player.ID] = player.Cards
		}
	}

	return visibleCards
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
