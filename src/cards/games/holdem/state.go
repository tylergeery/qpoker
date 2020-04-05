package holdem

import (
	"fmt"
	"qpoker/cards"
)

const (
	stateInit  = "Init"
	stateDeal  = "Deal"
	stateFlop  = "Flop"
	stateTurn  = "Turn"
	stateRiver = "River"
)

// HoldEm contains
type HoldEm struct {
	Board   []cards.Card
	Deck    cards.Deck
	Players []Player
	Pot     Pot
	State   string
	dealer  int
}

// NewHoldEm creates and returns the resources for a new HoldEm
func NewHoldEm(players []Player) *HoldEm {
	return &HoldEm{
		Deck:    cards.NewDeck(),
		Players: players,
		State:   stateInit,
	}
}

// Deal a new hand of holdem
func (h *HoldEm) Deal() *HoldEm {
	// shuffle cards
	h.Deck.Shuffle()

	// reset players cards
	for i := range h.Players {
		h.Players[i].Cards = make(cards.Card, 2)
	}

	// deal player cards
	for i := 0; i < 2; i++ {
		for j := 0; j < len(h.Players); j++ {
			h.Players[j].Cards = append(h.Players[j].Cards, h.Deck.GetCard())
		}
	}

	h.State = stateDeal

	return h
}

// Simulate game until completion
func (h *HoldEm) Simulate() error {
	for h.State != stateRiver {
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
		stateDeal: h.flop,
		stateFlop: h.turn,
		stateTurn: h.river,
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

	h.State = stateFlop
}

func (h *HoldEm) turn() {
	_ = h.Deck.GetCard() // Burn card

	h.addCardToBoard()
	h.State = stateTurn
}

func (h *HoldEm) river() {
	_ = h.Deck.GetCard() // Burn card

	h.addCardToBoard()
	h.State = stateRiver
}
