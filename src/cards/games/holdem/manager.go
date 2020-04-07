package holdem

import (
	"fmt"
	"qpoker/cards/games"
)

// GameOptions is a config object used to create a new GameManager
type GameOptions struct {
	Capacity    int   `json:"capacity"`
	BigBlind    int64 `json:"big_blind"`
	LittleBlind int64 `json:"little_blind"`
}

// GameManager holds game state and manages game flow
type GameManager struct {
	Table       *games.Table `json:"table"`
	State       *HoldEm      `json:"state"`
	Active      int          `json:"active"`
	Dealer      int          `json:"dealer"`
	ToPlay      int64        `json:"to_play"`
	BigBlind    int64        `json:"big_blind"`
	LittleBlind int64        `json:"little_blind"`
}

// NewGameManager returns a new GameManager
func NewGameManager(players []*games.Player, options GameOptions) (*GameManager, error) {
	state, err := NewHoldEm(players)
	if err != nil {
		return nil, err
	}

	gm := &GameManager{
		Table: &games.Table{
			Players:  players,
			Capacity: options.Capacity,
		},
		State:       state,
		Active:      0,
		Dealer:      0, // TODO: random?,
		ToPlay:      0,
		BigBlind:    options.BigBlind,
		LittleBlind: options.LittleBlind,
	}

	return gm, nil
}

func (g *GameManager) NextHand() {
	g.Dealer = g.nextPos(g.Dealer)

	// LittleBlind
	g.Active = g.nextPos(g.Dealer)
	g.playerBet(ActionNewBet(g.LittleBlind))

	// BigBlind
	g.Active = g.nextPos(g.Active)
	g.playerBet(ActionNewBet(g.BigBlind))

	g.ToPlay = g.BigBlind
	g.Active = g.nextPos(g.Active)
}

func (g *GameManager) nextPos(pos int) int {
	return (pos + 1) % len(g.State.Players)
}

func (g *GameManager) getRemainingPlayers() []*games.Player {
	remaining := []*games.Player{}

	for i := range g.State.Players {
		if g.State.Players[i].Active {
			remaining = append(remaining, g.State.Players[i])
		}
	}

	return remaining
}
func (g *GameManager) isComplete() bool {
	playersRemaining := g.getRemainingPlayers()

	if len(playersRemaining) <= 1 {
		return true
	}

	return false
}

func (g *GameManager) playerBet(action Action) error {
	return nil
}

func (g *GameManager) playerCall(action Action) error {
	return nil
}

func (g *GameManager) playerCheck(action Action) error {
	return nil
}

func (g *GameManager) playerFold(action Action) error {
	return nil
}

// PlayerAction performs an action for player
func (g *GameManager) PlayerAction(playerID int64, action Action) (complete bool, err error) {
	if g.isComplete() {
		err = fmt.Errorf("Game is already complete")
		return
	}

	if playerID != g.State.Players[g.Active].ID {
		err = fmt.Errorf("User (%d) must wait for player (%d) to act", playerID, g.State.Players[g.Active].ID)
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

	g.Active = g.nextPos(g.Active)

	complete = g.isComplete()

	return
}
