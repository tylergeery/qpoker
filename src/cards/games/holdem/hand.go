package holdem

import (
	"math"
	"qpoker/cards"
)

// Hand holds cards
type Hand struct {
	Cards []cards.Card
}

// GetScore return a numeric value for a hand
func (h Hand) GetScore() int64 {
	score := int64(0)

	for i, card := range h.Cards {
		score += int64(math.Pow10(2*i)) * int64(card.Value)
	}

	return score
}

// HandEvaluator takes a hand
type HandEvaluator struct {
	Hand   Hand
	aggMap map[int][]int
}

// NewHandEvaluator returns a new hand evaluator
func NewHandEvaluator(hand Hand) *HandEvaluator {
	counts := map[int]int{}
	aggMap := map[int][]int{}

	for _, card := range hand.Cards {
		if _, ok := counts[card.Value]; !ok {
			counts[card.Value] = 0
		}
		counts[card.Value] += 1
	}

	for val, count := range counts {
		if _, ok := aggMap[count]; !ok {
			aggMap[count] = []int{}
		}
		aggMap[count] = append(aggMap[count], val)
	}

	// TODO: ensure we have a sorted hand

	return &HandEvaluator{hand, aggMap}
}

// IsRoyalFlush checks for a royal flush
func (h HandEvaluator) IsRoyalFlush() (bool, Hand) {
	return true, h.Hand
}

// IsStraightFlush checks for a straight flush
func (h HandEvaluator) IsStraightFlush() (bool, Hand) {
	return true, h.Hand
}

// IsFourOfAKind checks for a four of a kind
func (h HandEvaluator) IsFourOfAKind() (bool, Hand) {
	return true, h.Hand
}

// IsFullHouse checks for a full house
func (h HandEvaluator) IsFullHouse() (bool, Hand) {
	return true, h.Hand
}

// IsFlush checks for a flush
func (h HandEvaluator) IsFlush() (bool, Hand) {
	return true, h.Hand
}

// IsStraight checks for a straight
func (h HandEvaluator) IsStraight() (bool, Hand) {
	return true, h.Hand
}

// IsThreeOfAKind checks for 3 of a kind
func (h HandEvaluator) IsThreeOfAKind() (bool, Hand) {
	return true, h.Hand
}

// IsTwoPair checks for two pair
func (h HandEvaluator) IsTwoPair() (bool, Hand) {
	return true, h.Hand
}

// IsPair checks for a pair
func (h HandEvaluator) IsPair() (bool, Hand) {
	return true, h.Hand
}

// IsHighCard checks for a high card (always true)
func (h HandEvaluator) IsHighCard() (bool, Hand) {
	return true, h.Hand
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

	for _, priorityFunction := range handPriorities {
		if match, bestHand := priorityFunction(); match {
			return bestHand.GetScore()
		}
	}

	return int64(0)
}
