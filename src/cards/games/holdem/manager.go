package holdem

import (
	"fmt"
	"qpoker/utils"
)

const (
	// MaxPlayerCount is the max amount of players for HoldEm
	MaxPlayerCount = 12
)

// GameOptions is a config object used to create a new GameManager
type GameOptions struct {
	Capacity int   `json:"capacity"`
	BigBlind int64 `json:"big_blind"`
}

// GameManager holds game state and manages game flow
type GameManager struct {
	State    *HoldEm `json:"state"`
	Pot      *Pot    `json:"pot"`
	BigBlind int64   `json:"big_blind"`
}

// NewGameManager returns a new GameManager
func NewGameManager(players []*Player, options GameOptions) (*GameManager, error) {
	if len(players) <= 1 || len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	if options.Capacity > 0 && len(players) > options.Capacity {
		return nil, fmt.Errorf("Player count (%d) is greater than capacity (%d)", len(players), options.Capacity)
	}
	table := NewTable(options.Capacity, players)
	gm := &GameManager{
		State:    NewHoldEm(table),
		BigBlind: options.BigBlind,
	}

	return gm, nil
}

// NextHand moves game manager on to next hand
func (g *GameManager) NextHand() error {
	err := g.State.Deal()
	if err != nil {
		return err
	}

	g.Pot = NewPot(g.State.Table.GetActivePlayers())

	// little
	g.State.Table.GetActivePlayer().LittleBlind = true
	g.playerBet(NewActionBet(g.BigBlind / 2))
	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	// big blind
	g.State.Table.GetActivePlayer().BigBlind = true
	g.playerBet(NewActionBet(g.BigBlind))
	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	return nil
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
		g.Pot.GetPayouts(g.State.GetWinningIDs())
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
