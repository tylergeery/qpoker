package holdem

import (
	"math"
	"qpoker/cards"
	"sort"
)

// Hand holds cards
type Hand struct {
	Cards []cards.Card
}

// FormBestHand returns the bests possible hand
func FormBestHand(required []cards.Card, remaining []cards.Card) Hand {
	return Hand{Cards: append(required, remaining...)[:5]}
}

// GetScore return a numeric value for a hand
func (h Hand) GetScore(handScore int) int64 {
	score := int64(0)

	for i, card := range h.Cards {
		score += int64(math.Pow10(2*i)) * int64(card.Value)
	}

	score += int64(math.Pow10(2*len(h.Cards))) * int64(handScore)

	return score
}

// HandEvaluator takes a hand
type HandEvaluator struct {
	Hand        Hand
	aggValueMap map[int][]int
	aggSuitMap  map[int][]byte
}

// NewHandEvaluator returns a new hand evaluator
func NewHandEvaluator(hand Hand) *HandEvaluator {
	valueCounts := map[int]int{}
	suitCounts := map[byte]int{}

	aggValueMap := map[int][]int{}
	aggSuitMap := map[int][]byte{}

	for _, card := range hand.Cards {
		if _, ok := valueCounts[card.Value]; !ok {
			valueCounts[card.Value] = 0
		}
		if _, ok := suitCounts[card.Suit]; !ok {
			suitCounts[card.Suit] = 0
		}
		suitCounts[card.Suit]++
		valueCounts[card.Value]++
	}

	for val, count := range valueCounts {
		if _, ok := aggValueMap[count]; !ok {
			aggValueMap[count] = []int{}
		}
		aggValueMap[count] = append(aggValueMap[count], val)
	}

	for suit, count := range suitCounts {
		if count > 5 {
			count = 5
		}

		if _, ok := aggSuitMap[count]; !ok {
			aggSuitMap[count] = []byte{}
		}
		aggSuitMap[count] = append(aggSuitMap[count], suit)
	}

	// Sort values for ties
	for count := range aggValueMap {
		sort.Ints(aggValueMap[count])
	}

	// Evaluator always expects hand to be sorted
	sort.Slice(
		hand.Cards,
		func(i, j int) bool {
			if hand.Cards[i].Value == 1 && hand.Cards[j].Value != 1 {
				return true
			}

			return hand.Cards[i].Value > hand.Cards[j].Value
		},
	)

	return &HandEvaluator{hand, aggValueMap, aggSuitMap}
}

func (h HandEvaluator) getCardsBySuit(suit byte) []cards.Card {
	resp := []cards.Card{}

	for _, card := range h.Hand.Cards {
		if card.Suit == suit {
			resp = append(resp, card)
		}
	}

	return resp
}

func (h HandEvaluator) getCardsWithoutSuit(suit byte) []cards.Card {
	resp := []cards.Card{}

	for _, card := range h.Hand.Cards {
		if card.Suit != suit {
			resp = append(resp, card)
		}
	}

	return resp
}

func (h HandEvaluator) getCardsByValue(value int) []cards.Card {
	resp := []cards.Card{}

	for _, card := range h.Hand.Cards {
		if card.Value == value {
			resp = append(resp, card)
		}
	}

	return resp
}

func (h HandEvaluator) getCardsWithoutValue(value int) []cards.Card {
	resp := []cards.Card{}

	for _, card := range h.Hand.Cards {
		if card.Value != value {
			resp = append(resp, card)
		}
	}

	return resp
}

// IsRoyalFlush checks for a royal flush
func (h HandEvaluator) IsRoyalFlush() (bool, Hand) {
	valid, bestHand := h.IsStraightFlush()
	if !valid {
		return false, h.Hand
	}

	if bestHand.Cards[0].Value != 1 {
		return false, h.Hand
	}

	return true, bestHand
}

// IsStraightFlush checks for a straight flush
func (h HandEvaluator) IsStraightFlush() (bool, Hand) {
	suits, ok := h.aggSuitMap[5]
	if !ok {
		return false, h.Hand
	}

	suited := h.getCardsBySuit(suits[0])
	straight := []cards.Card{suited[0]}

	for i := 1; i < len(suited); i++ {
		prev := straight[len(straight)-1]
		if prev.MatchesValue(suited[i]) {
			continue
		}

		if !prev.Connects(suited[i]) {
			straight = []cards.Card{}
		}

		straight = append(straight, suited[i])
		if len(straight) == 5 {
			return true, Hand{straight}
		}
	}

	return false, h.Hand
}

// IsFourOfAKind checks for a four of a kind
func (h HandEvaluator) IsFourOfAKind() (bool, Hand) {
	vals, ok := h.aggValueMap[4]
	if !ok {
		return false, h.Hand
	}

	ofKindCards := h.getCardsByValue(vals[0])
	nextBestCards := h.getCardsWithoutValue(vals[0])

	return true, FormBestHand(ofKindCards, nextBestCards)
}

// IsFullHouse checks for a full house
func (h HandEvaluator) IsFullHouse() (bool, Hand) {
	threeVal, threeOk := h.aggValueMap[3]
	twoVal, twoOk := h.aggValueMap[2]

	if !threeOk || !twoOk {
		return false, h.Hand
	}

	threeOfKindCards := h.getCardsByValue(threeVal[0])
	twoOfKindCards := h.getCardsByValue(twoVal[0])

	return true, FormBestHand(threeOfKindCards, twoOfKindCards)
}

// IsFlush checks for a flush
func (h HandEvaluator) IsFlush() (bool, Hand) {
	suits, ok := h.aggSuitMap[5]
	if !ok {
		return false, h.Hand
	}

	ofKindCards := h.getCardsBySuit(suits[0])
	nextBestCards := h.getCardsWithoutSuit(suits[0])

	return true, FormBestHand(ofKindCards, nextBestCards)
}

// IsStraight checks for a straight
func (h HandEvaluator) IsStraight() (bool, Hand) {
	allCards := h.Hand.Cards
	straight := []cards.Card{allCards[0]}

	for i := 1; i < len(allCards); i++ {
		prev := straight[len(straight)-1]
		if prev.MatchesValue(allCards[i]) {
			continue
		}

		if !prev.Connects(allCards[i]) {
			straight = []cards.Card{}
		}

		straight = append(straight, allCards[i])
		if len(straight) == 5 {
			return true, Hand{straight}
		}
	}

	return false, h.Hand
}

// IsThreeOfAKind checks for 3 of a kind
func (h HandEvaluator) IsThreeOfAKind() (bool, Hand) {
	threeVal, threeOk := h.aggValueMap[3]

	if !threeOk {
		return false, h.Hand
	}

	threeOfKindCards := h.getCardsByValue(threeVal[0])
	remCards := h.getCardsWithoutValue(threeVal[0])

	return true, FormBestHand(threeOfKindCards, remCards)
}

// IsTwoPair checks for two pair
func (h HandEvaluator) IsTwoPair() (bool, Hand) {
	twoVals, twoOk := h.aggValueMap[3]

	if !twoOk || len(twoVals) <= 1 {
		return false, h.Hand
	}

	firstPairCards := h.getCardsByValue(twoVals[0])
	secondPairCards := h.getCardsWithoutValue(twoVals[1])
	remCards := h.getCardsWithoutValue(twoVals[0])
	nextHand := HandEvaluator{Hand: Hand{Cards: remCards}}
	remCards = nextHand.getCardsWithoutValue(twoVals[1])

	return true, FormBestHand(append(firstPairCards, secondPairCards...), remCards)
}

// IsPair checks for a pair
func (h HandEvaluator) IsPair() (bool, Hand) {
	twoVals, twoOk := h.aggValueMap[3]

	if !twoOk {
		return false, h.Hand
	}

	twoOfKindCards := h.getCardsByValue(twoVals[0])
	remCards := h.getCardsWithoutValue(twoVals[0])

	return true, FormBestHand(twoOfKindCards, remCards)
}

// IsHighCard checks for a high card (always true)
func (h HandEvaluator) IsHighCard() (bool, Hand) {
	return true, FormBestHand(h.Hand.Cards[:5], nil)
}

// Evaluate returns a numeric score for a hand (higher the better)
func Evaluate(hand Hand) int64 {
	handEvaluator := NewHandEvaluator(hand)
	handPriorities := []func() (bool, Hand){
		handEvaluator.IsRoyalFlush,
		handEvaluator.IsStraightFlush,
		handEvaluator.IsFourOfAKind,
		handEvaluator.IsFullHouse,
		handEvaluator.IsFlush,
		handEvaluator.IsStraight,
		handEvaluator.IsThreeOfAKind,
		handEvaluator.IsTwoPair,
		handEvaluator.IsPair,
		handEvaluator.IsHighCard,
	}

	for i, priorityFunction := range handPriorities {
		if match, bestHand := priorityFunction(); match {
			return bestHand.GetScore(len(handPriorities) - i)
		}
	}

	return int64(0)
}
