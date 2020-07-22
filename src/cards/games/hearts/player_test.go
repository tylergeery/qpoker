package hearts

import (
	"qpoker/cards"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlayerAddCards(t *testing.T) {
	type TestCase struct {
		initial  []cards.Card
		add      []cards.Card
		expected []cards.Card
	}

	testCases := []TestCase{
		TestCase{
			initial: []cards.Card{
				cards.NewCard(10, 'C'), cards.NewCard(1, 'D'), cards.NewCard(1, 'C'),
				cards.NewCard(6, 'C'), cards.NewCard(2, 'S'), cards.NewCard(1, 'S'),
				cards.NewCard(5, 'H'), cards.NewCard(4, 'S'), cards.NewCard(1, 'H'),
				cards.NewCard(10, 'H'),
			},
			add: []cards.Card{},
			expected: []cards.Card{
				cards.NewCard(6, 'C'), cards.NewCard(10, 'C'), cards.NewCard(1, 'C'),
				cards.NewCard(1, 'D'), cards.NewCard(2, 'S'), cards.NewCard(4, 'S'),
				cards.NewCard(1, 'S'), cards.NewCard(5, 'H'), cards.NewCard(10, 'H'),
				cards.NewCard(1, 'H'),
			},
		},
		TestCase{
			initial: []cards.Card{
				cards.NewCard(10, 'C'), cards.NewCard(11, 'C'), cards.NewCard(1, 'C'),
				cards.NewCard(6, 'C'), cards.NewCard(2, 'S'), cards.NewCard(1, 'S'),
				cards.NewCard(5, 'H'), cards.NewCard(4, 'S'), cards.NewCard(1, 'H'),
				cards.NewCard(10, 'H'),
			},
			add: []cards.Card{
				cards.NewCard(5, 'D'), cards.NewCard(12, 'D'), cards.NewCard(10, 'D'),
			},
			expected: []cards.Card{
				cards.NewCard(6, 'C'), cards.NewCard(10, 'C'), cards.NewCard(11, 'C'),
				cards.NewCard(1, 'C'), cards.NewCard(5, 'D'), cards.NewCard(10, 'D'),
				cards.NewCard(12, 'D'), cards.NewCard(2, 'S'), cards.NewCard(4, 'S'),
				cards.NewCard(1, 'S'), cards.NewCard(5, 'H'), cards.NewCard(10, 'H'),
				cards.NewCard(1, 'H'),
			},
		},
	}

	for _, c := range testCases {
		p := &Player{Cards: c.initial}
		p.AddCards(c.add)

		for i := range c.expected {
			assert.Equal(t, c.expected[i].ToString(), p.Cards[i].ToString())
		}
	}
}
