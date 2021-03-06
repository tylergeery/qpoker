package holdem

import (
	"fmt"
	"qpoker/cards"
	"sort"
)

const (
	// StateInit initial game state
	StateInit = "Init"
	// StateDeal game state after all players have cards
	StateDeal = "Deal"
	// StateFlop game state after initial 3 cards flipped
	StateFlop = "Flop"
	// StateTurn game state after 4th card turned
	StateTurn = "Turn"
	// StateRiver game state after last card turned
	StateRiver = "River"
)

// HoldEm contains game state
type HoldEm struct {
	Board []cards.Card `json:"board"`
	Deck  *cards.Deck  `json:"-"`
	Table *Table       `json:"table"`
	State string       `json:"state"`
}

// NewHoldEm creates and returns the resources for a new HoldEm
func NewHoldEm(table *Table) *HoldEm {
	return &HoldEm{
		Deck:  cards.NewDeck(),
		Table: table,
		State: StateInit,
	}
}

// Deal a new hand of holdem
func (h *HoldEm) Deal() error {
	// shuffle cards
	h.Deck.Shuffle()

	// reset table
	err := h.Table.NextHand()
	if err != nil {
		return err
	}

	// reset board
	h.Board = []cards.Card{}

	// reset players cards
	players := h.Table.GetActivePlayers()
	for i := range players {
		players[i].Cards = []cards.Card{}
		players[i].CardsVisible = false
	}

	// deal player cards
	for i := 0; i < 2; i++ {
		for j := 0; j < len(players); j++ {
			card, err := h.Deck.GetCard()
			if err != nil {
				return fmt.Errorf("unexpected GetCard deal error: %s", err)
			}

			players[j].Cards = append(players[j].Cards, card)
		}
	}

	h.State = StateDeal

	return nil
}

// Simulate game until completion
func (h *HoldEm) Simulate() error {
	for h.State != StateRiver {
		err := h.Advance()
		if err != nil {
			return err
		}
	}

	return nil
}

// Advance game to next state
func (h *HoldEm) Advance() error {
	advanceStateMap := map[string]func(){
		StateDeal: h.flop,
		StateFlop: h.turn,
		StateTurn: h.river,
	}

	advance, ok := advanceStateMap[h.State]
	if !ok {
		return fmt.Errorf("Game cannot advance from state: %s", h.State)
	}

	advance()

	return nil
}

func toTieredWinners(orderedPlayerIDs []int64, scores map[int64]int64) [][]int64 {
	tiered := [][]int64{}
	tier := []int64{}

	for _, playerID := range orderedPlayerIDs {
		// players tied and belong in same tier
		if len(tier) == 0 || scores[tier[0]] == scores[playerID] {
			tier = append(tier, playerID)
			continue
		}

		// start new tier
		tiered = append(tiered, tier)
		tier = []int64{playerID}
	}

	tiered = append(tiered, tier)

	return tiered
}

// GetWinningIDs returns the hand winners in order
func (h *HoldEm) GetWinningIDs() [][]int64 {
	players := h.Table.GetActivePlayers()
	scores := map[int64]int64{}
	playerIDs := []int64{}

	if len(players) == 1 {
		return [][]int64{[]int64{players[0].ID}}
	}

	for i := range players {
		handCards := []cards.Card{}
		handCards = append(handCards, h.Board...)
		handCards = append(handCards, players[i].Cards...)

		scores[players[i].ID] = Evaluate(Hand{handCards})
		playerIDs = append(playerIDs, players[i].ID)
	}

	sort.Slice(playerIDs, func(i, j int) bool {
		return scores[playerIDs[i]] >= scores[playerIDs[j]]
	})

	fmt.Printf("WinningIDs: %+v, scores: %+v\n", playerIDs, scores)

	return toTieredWinners(playerIDs, scores)
}

func (h *HoldEm) addCardToBoard() {
	card, err := h.Deck.GetCard()
	if err != nil {
		fmt.Printf("Unexpected GetCard error: %s", err)
		return
	}

	h.Board = append(h.Board, card)
}
func (h *HoldEm) flop() {
	_, _ = h.Deck.GetCard() // Burn card

	for i := 0; i < 3; i++ {
		h.addCardToBoard()
	}

	h.State = StateFlop
}

func (h *HoldEm) turn() {
	_, _ = h.Deck.GetCard() // Burn card

	h.addCardToBoard()
	h.State = StateTurn
}

func (h *HoldEm) river() {
	_, _ = h.Deck.GetCard() // Burn card

	h.addCardToBoard()
	h.State = StateRiver
}
