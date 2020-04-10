package holdem

import (
	"qpoker/cards"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	type TestCase struct {
		hand          Hand
		expectedScore int64
	}
	cases := []TestCase{
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(9, cards.SuitHearts),
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(13, cards.SuitClubs),
					cards.NewCard(11, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(5, cards.SuitSpades),
				},
			},
			expectedScore: int64(11311100905),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(13, cards.SuitClubs),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(4, cards.SuitClubs),
					cards.NewCard(9, cards.SuitClubs),
				},
			},
			expectedScore: int64(11413111009),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(3, cards.SuitClubs),
					cards.NewCard(4, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(13, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(11413100804),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(6, cards.SuitClubs),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(21414100806),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(11, cards.SuitHearts),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(2, cards.SuitClubs),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(30303020211),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(11, cards.SuitHearts),
					cards.NewCard(11, cards.SuitHearts),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(41111111410),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(9, cards.SuitClubs),
					cards.NewCard(12, cards.SuitHearts),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(51211100908),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(9, cards.SuitClubs),
					cards.NewCard(12, cards.SuitHearts),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(13, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(8, cards.SuitSpades),
				},
			},
			expectedScore: int64(51413121110),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(4, cards.SuitClubs),
					cards.NewCard(5, cards.SuitHearts),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(3, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(2, cards.SuitSpades),
				},
			},
			expectedScore: int64(50504030214),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(9, cards.SuitClubs),
					cards.NewCard(12, cards.SuitHearts),
					cards.NewCard(11, cards.SuitDiamonds),
					cards.NewCard(6, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitDiamonds),
					cards.NewCard(2, cards.SuitDiamonds),
				},
			},
			expectedScore: int64(61411100602),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(9, cards.SuitDiamonds),
					cards.NewCard(12, cards.SuitDiamonds),
					cards.NewCard(4, cards.SuitDiamonds),
					cards.NewCard(13, cards.SuitDiamonds),
					cards.NewCard(10, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitDiamonds),
					cards.NewCard(8, cards.SuitDiamonds),
				},
			},
			expectedScore: int64(61413121009),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(2, cards.SuitClubs),
					cards.NewCard(11, cards.SuitDiamonds),
					cards.NewCard(11, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(5, cards.SuitSpades),
				},
			},
			expectedScore: int64(70202021111),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(1, cards.SuitDiamonds),
					cards.NewCard(1, cards.SuitClubs),
					cards.NewCard(1, cards.SuitSpades),
					cards.NewCard(2, cards.SuitDiamonds),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(5, cards.SuitSpades),
				},
			},
			expectedScore: int64(81414141405),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(1, cards.SuitHearts),
					cards.NewCard(2, cards.SuitHearts),
					cards.NewCard(3, cards.SuitHearts),
					cards.NewCard(4, cards.SuitHearts),
					cards.NewCard(5, cards.SuitHearts),
					cards.NewCard(6, cards.SuitSpades),
					cards.NewCard(5, cards.SuitSpades),
				},
			},
			expectedScore: int64(90504030214),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(10, cards.SuitClubs),
					cards.NewCard(13, cards.SuitClubs),
					cards.NewCard(12, cards.SuitClubs),
					cards.NewCard(9, cards.SuitClubs),
					cards.NewCard(12, cards.SuitDiamonds),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(1, cards.SuitHearts),
				},
			},
			expectedScore: int64(91312111009),
		},
		TestCase{
			hand: Hand{
				Cards: []cards.Card{
					cards.NewCard(10, cards.SuitClubs),
					cards.NewCard(13, cards.SuitClubs),
					cards.NewCard(12, cards.SuitClubs),
					cards.NewCard(12, cards.SuitSpades),
					cards.NewCard(12, cards.SuitDiamonds),
					cards.NewCard(11, cards.SuitClubs),
					cards.NewCard(1, cards.SuitClubs),
				},
			},
			expectedScore: int64(101413121110),
		},
	}

	for _, c := range cases {
		assert.Equal(t, c.expectedScore, Evaluate(c.hand))
	}
}
