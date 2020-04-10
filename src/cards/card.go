package cards

import (
	"fmt"
	"strconv"
)

const (
	// SuitClubs is byte character for clubs
	SuitClubs = 'C'
	// SuitDiamonds is byte character for diamonds
	SuitDiamonds = 'D'
	// SuitHearts is byte character for hearts
	SuitHearts = 'H'
	// SuitSpades is byte character for spades
	SuitSpades = 'S'
)

// Card is a single card object
type Card struct {
	Value int
	Suit  byte
	Char  byte
}

// NewCard returns a card from card string value
func NewCard(value int, suit byte) Card {
	cardmap := map[int]byte{
		1:  'A',
		10: 'T',
		11: 'J',
		12: 'Q',
		13: 'K',
	}

	if char, ok := cardmap[value]; ok {
		return Card{value, suit, char}
	}

	return Card{value, suit, strconv.Itoa(value)[0]}
}

// ToString gets string representation of a card
func (c Card) ToString() string {
	return fmt.Sprintf("%c%c", c.Char, c.Suit)
}

// Connects determines whether a card connects to form a run with another
func (c Card) Connects(next Card) bool {
	if c.Value == 1 {
		return next.Value == 2 || next.Value == 13
	}

	if next.Value == 1 {
		return c.Value == 2 || c.Value == 13
	}

	return c.Value == (next.Value+1) || c.Value == (next.Value-1)
}

// MatchesValue is whether a card matches another by value
func (c Card) MatchesValue(next Card) bool {
	return c.Value == next.Value
}

// MatchesSuit is whether a card matches another by suit
func (c Card) MatchesSuit(next Card) bool {
	return c.Suit == next.Suit
}
