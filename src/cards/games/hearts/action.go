package holdem

import (
	"qpoker/cards"
)

const (
	// ActionPlay is a play action
	ActionPlay = "play"
)

// Action holds the info regarding a holdem action
type Action struct {
	Name string
	Card cards.Card
}

func getCardsForAction(cardStrings []string) []cards.Card {
	actionCards := []cards.Card{}

	for _, c := range cardStrings {
		actionCards = append(actionCards, cards.CardFromString(c))
	}

	return actionCards
}

// NewActionPlay returns a new play action
func NewActionPlay(c string) Action {
	return Action{ActionPlay, cards.CardFromString(c)}
}
