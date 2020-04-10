package cards

import (
	"fmt"
	"math/rand"
)

// Deck represent a single deck of cards
type Deck struct {
	Cards []Card
	index int
}

// NewDeck returns a new deck of cards for use
func NewDeck() *Deck {
	suits := getSuits()
	cards := getCards()
	deck := &Deck{
		Cards: []Card{},
		index: 0,
	}

	for _, suit := range suits {
		for _, card := range cards {
			deck.Cards = append(deck.Cards, NewCard(card, suit))
		}
	}

	return deck
}

func getSuits() []byte {
	return []byte{SuitClubs, SuitDiamonds, SuitHearts, SuitSpades}
}

func getCards() []int {
	cards := []int{}

	for i := 1; i <= 13; i++ {
		cards = append(cards, i)
	}

	return cards
}

// Shuffle returns the deck shuffled
func (d *Deck) Shuffle() *Deck {
	var next int

	for i := range d.Cards {
		for true {
			next = rand.Intn(len(d.Cards))
			if next != i {
				break
			}
		}

		d.Cards[i], d.Cards[next] = d.Cards[next], d.Cards[i]
	}

	d.index = 0

	return d
}

// GetCard returns the next card in the deck, or an error
func (d *Deck) GetCard() (card Card, err error) {
	if d.index == len(d.Cards) {
		err = fmt.Errorf("No more cards in deck")
		return
	}

	card = d.Cards[d.index]
	d.index++

	return
}
