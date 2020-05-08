package cards

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeck(t *testing.T) {
	// Validate Deck
	deck := NewDeck()

	cardsBySuit := map[byte][]int{
		SuitClubs:    make([]int, 13),
		SuitDiamonds: make([]int, 13),
		SuitHearts:   make([]int, 13),
		SuitSpades:   make([]int, 13),
	}
	for _, card := range deck.Cards {
		cardsBySuit[card.Suit][card.Value-1] = card.Value
	}

	// Ensure a full deck
	for _, cards := range cardsBySuit {
		for i, val := range cards {
			assert.Equal(t, i, val-1)
		}
	}
}

func TestDealAndShuffle(t *testing.T) {
	deck := NewDeck()
	deck.Shuffle()
	assert.Equal(t, 0, deck.index)

	first, _ := deck.GetCard()
	second, _ := deck.GetCard()
	third, _ := deck.GetCard()

	// Shuffle
	deck.Shuffle()
	assert.Equal(t, 0, deck.index)

	first1, _ := deck.GetCard()
	second2, _ := deck.GetCard()
	third3, _ := deck.GetCard()

	assert.NotEqual(t, first, first1)
	assert.NotEqual(t, second, second2)
	assert.NotEqual(t, third, third3)
}

func TestCardStringValues(t *testing.T) {
	deck := NewDeck()

	for i := 0; i < 5; i++ {
		deck.Shuffle()
		cardMap := map[string]bool{}

		for j := 0; j < 52; j++ {
			card, err := deck.GetCard()
			assert.NoError(t, err)
			_, ok := cardMap[card.ToString()]
			assert.False(t, ok, fmt.Sprintf("Found multiple card values in deck: %s, %+v", card.ToString(), card))
			cardMap[card.ToString()] = true
		}
	}
}

func TestGetCardError(t *testing.T) {
	deck := NewDeck()

	for i := 0; i < 5; i++ {
		deck.Shuffle()
		for j := 0; j < 52; j++ {
			_, err := deck.GetCard()
			assert.NoError(t, err)
		}

		_, err := deck.GetCard()
		assert.Error(t, err)
	}
}
