package cards

import (
	"fmt"
	"math/rand"
	"strconv"
)

const (
	SUIT_CLUBS    = "C"
	SUIT_DIAMONDS = "D"
	SUIT_HEARTS   = "H"
	SUIT_SPADES   = "S"
)

// Card is a single card object
type Card struct {
	value string
	suit  string
}

// Deck represent a single deck of cards
type Deck struct {
	Cards []Card
	index int
}

func getSuits() []string {
	return []string{SUIT_CLUBS, SUIT_DIAMONDS, SUIT_HEARTS, SUIT_SPADES}
}

func getCards() []string {
	cards := []string{}
	cardmap := map[int]string{
		1:  "A",
		11: "J",
		12: "Q",
		13: "K",
	}

	for i := 1; i < 13; i++ {
		card := strconv.Itoa(i)
		if _card, ok := cardmap[i]; ok {
			card = _card
		}

		cards = append(cards, card)
	}

	return cards
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
			deck.Cards = append(deck.Cards, Card{card, suit})
		}
	}

	return deck
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
