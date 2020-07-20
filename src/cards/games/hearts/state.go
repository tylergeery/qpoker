package holdem

import (
	"fmt"
	"qpoker/cards"
)

const (
	// PassLeft means pass left
	PassLeft = 'l'
	// PassRight means pass right
	PassRight = 'r'
	// PassAcross means pass across
	PassAcross = 'a'
	// PassNone means no pass
	PassNone = 'n'

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

	passes map[int64][]cards.Card
}

// NewHearts creates and returns the resources for a new Hearts game
func NewHearts(table *Table, state string) *Hearts {
	scores := map[int64]int{}
	for _, player := range table.GetAllPlayers() {
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
	// reset table
	h.Table.NextHand()

	// reset board
	h.Board = map[int64]cards.Card{}
	h.passes = map[int64][]cards.Card{}

	// shuffle cards
	h.Deck.Shuffle()

	// deal player cards
	next := h.Table.ActiveIndex
	players := h.Table.GetActivePlayers()

	for i := 0; i < 52; i++ {
		card, err := h.Deck.GetCard()
		if err != nil {
			return fmt.Errorf("unexpected GetCard deal error hearts: %s", err)
		}

		players[next].Cards = append(players[next].Cards, card)
		next = h.Table.NextPos(next)
	}

	h.State = StatePassing

	return nil
}

func (h *Hearts) clearBoard() {
	h.Board = map[int64]cards.Card{}
}

func (h *Hearts) playerPlay(playerID int64, card cards.Card) {
	h.Board[playerID] = card

	// Remove card from players
	player := h.Table.GetPlayerByID(playerID)
	for i := range player.Cards {
		if player.Cards[i].ToString() == card.ToString() {
			player.Cards = append(player.Cards[:i], player.Cards[i+1:]...)
			break
		}
	}
}

func (h *Hearts) addPass(playerID int64, cards []cards.Card) {
	h.passes[playerID] = cards
}
