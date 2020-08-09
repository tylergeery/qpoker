package cards

import (
	"sort"
)

func sortSuited(cards []Card, suitOrdering []byte, acesHigh bool) {
	suitMap := map[byte]int{}
	for i, b := range suitOrdering {
		suitMap[b] = i
	}

	cardValue := func(card Card) int {
		if card.Value == 1 {
			return 14
		}

		return card.Value
	}

	sort.Slice(
		cards,
		func(i, j int) bool {
			cardA, cardB := cards[i], cards[j]

			if cardA.Suit != cardB.Suit {
				return suitMap[cardA.Suit] < suitMap[cardB.Suit]
			}

			return cardValue(cardA) < cardValue(cardB)
		},
	)
}

// SortSuitedAcesHigh sorts cards suit with Aces as high value
func SortSuitedAcesHigh(cards []Card) {
	sortSuited(cards, []byte{SuitClubs, SuitDiamonds, SuitSpades, SuitHearts}, true)
}
