package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardsToAndFromString(t *testing.T) {
	// Given
	type TestCase struct {
		card     Card
		expected string
	}

	testCases := []TestCase{
		TestCase{NewCard(1, 'C'), "AC"},
		TestCase{NewCard(2, 'D'), "2D"},
		TestCase{NewCard(3, 'S'), "3S"},
		TestCase{NewCard(4, 'H'), "4H"},
		TestCase{NewCard(5, 'D'), "5D"},
		TestCase{NewCard(6, 'S'), "6S"},
		TestCase{NewCard(7, 'D'), "7D"},
		TestCase{NewCard(8, 'H'), "8H"},
		TestCase{NewCard(9, 'C'), "9C"},
		TestCase{NewCard(10, 'S'), "TS"},
		TestCase{NewCard(11, 'D'), "JD"},
		TestCase{NewCard(12, 'H'), "QH"},
		TestCase{NewCard(13, 'C'), "KC"},
	}

	for _, tc := range testCases {
		s := tc.card.ToString()
		assert.Equal(t, tc.expected, s)
		card := CardFromString(s)
		assert.Equal(t, tc.card.Char, card.Char)
		assert.Equal(t, tc.card.Suit, card.Suit)
		assert.Equal(t, tc.card.Value, card.Value)
	}
}
