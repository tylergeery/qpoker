package holdem

import (
	"fmt"
	"qpoker/cards"
)

const (
	// StatePassing game state when all players are selecting cards to pass
	StatePassing = "Passing"
	// StateActive game state after game is active
	StateActive = "Flop"
)

// Hearts contains game state
type Hearts struct {
	Board  map[int64]cards.Card `json:"board"`
	Piles  map[int64]cards.Card `json:"piles"`
	Deck   *cards.Deck          `json:"-"`
	Table  *Table               `json:"table"`
	State  string               `json:"state"`
	Scores map[int64]int        `json:"scores"`
}

// NewHearts creates and returns the resources for a new Hearts game
func NewHearts(table *Table, state string) *Hearts {
	scores := map[int64]int{}
	for _, player := range table.GetPlayers() {
		scores[player.ID] = 0
	}

	return &Hearts{
		Deck:   cards.NewDeck(),
		Table:  table,
		State:  state,
		Scores: scores,
	}
}

// Deal a new hand of hearts
func (h *Hearts) Deal() error {
	// shuffle cards
	h.Deck.Shuffle()

	// reset table
	err := h.Table.NextHand()
	if err != nil {
		return err
	}

	// reset board
	h.Board = map[int64]cards.Card{}

	// reset players cards
	players := h.Table.GetActivePlayers()
	for i := range players {
		players[i].Cards = []cards.Card{}
		players[i].Pile = []cards.Card{}
	}

	// deal player cards
	next := h.Table.activeIndex
	for i := 0; i < 52; i++ {
		card, err := h.Deck.GetCard()
		if err != nil {
			return fmt.Errorf("unexpected GetCard deal error hearts: %s", err)
		}

		players[next].Cards = append(players[next].Cards, card)
		next = h.Table.nextPos(next)
	}

	h.State = StatePassing

	return nil
}
