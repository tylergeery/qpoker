package holdem

import (
	"fmt"
	"qpoker/cards"
	"qpoker/cards/games"
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
	Board   []cards.Card    `json:"board"`
	Deck    *cards.Deck     `json:"-"`
	Players []*games.Player `json:"-"`
	Pot     Pot             `json:"pot"`
	State   string          `json:"state"`
}

const (
	// MaxPlayerCount is the max amount of players for HoldEm
	MaxPlayerCount = 12
)

// NewHoldEm creates and returns the resources for a new HoldEm
func NewHoldEm(players []*games.Player) (*HoldEm, error) {
	if len(players) <= 1 || len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	return &HoldEm{
		Deck:    cards.NewDeck(),
		Players: players,
		State:   StateInit,
	}, nil
}

func (h *HoldEm) nextPos(n int) int {
	return (n + 1) % len(h.Players)
}

// Deal a new hand of holdem
func (h *HoldEm) Deal() *HoldEm {
	// shuffle cards
	h.Deck.Shuffle()

	// reset players cards
	for i := range h.Players {
		h.Players[i].Cards = []cards.Card{}
	}

	// deal player cards
	for i := 0; i < 2; i++ {
		for j := 0; j < len(h.Players); j++ {
			card, err := h.Deck.GetCard()
			if err != nil {
				fmt.Printf("Unexpected GetCard deal error: %s\n", err)
				return h
			}

			h.Players[j].Cards = append(h.Players[j].Cards, card)
		}
	}

	h.State = StateDeal

	return h
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
