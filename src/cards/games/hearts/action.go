package holdem

import (
	"qpoker/cards"
)

const (
	// ActionPlay is a play action
	ActionPlay = "bet"
	// ActionPass is a pass action
	ActionPass = "pass"
)

// Action holds the info regarding a holdem action
type Action struct {
	Name  string       `json:"name"`
	Cards []cards.Card `json:"cards"`
}

func getCardsForAction(cardStrings []string) []cards.Card {
	actionCards := []cards.Card{}

	for _, c := range cardStrings {
		actionCards = append(actionCards, cards.CardFromString(c))
	}

	return actionCards
}

// NewActionPlay returns a new play action
func NewActionPlay(cardStrings []string) Action {
	return Action{ActionPlay, getCardsForAction(cardStrings)}
}

// NewActionPass returns a new pass action
func NewActionPass(cardStrings []string) Action {
	return Action{ActionPass, getCardsForAction(cardStrings)}
}
