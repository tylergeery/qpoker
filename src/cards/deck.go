package cards

import (
	"fmt"
	"math/rand"
	"strconv"
)

const (
	SUIT_CLUBS    = 'C'
	SUIT_DIAMONDS = 'D'
	SUIT_HEARTS   = 'H'
	SUIT_SPADES   = 'S'
)

// Card is a single card object
type Card struct {
	Value int
	Suit  byte
	Char  byte
}

// ToString gets string representation of a card
func (c Card) ToString() string {
	return fmt.Sprintf("%c%c", c.Char, c.Suit)
}

// NewCard returns a card from card string value
func NewCard(value int, suit byte) Card {
	cardmap := map[int]byte{
		1:  'A',
		10: 'T',
		11: 'J',
		12: 'Q',
		13: 'K',
	}

	if char, ok := cardmap[value]; ok {
		return Card{value, suit, char}
	}

	return Card{value, suit, strconv.Itoa(value)[0]}
}

// Deck represent a single deck of cards
type Deck struct {
	Cards []Card
	index int
}

func getSuits() []byte {
	return []byte{SUIT_CLUBS, SUIT_DIAMONDS, SUIT_HEARTS, SUIT_SPADES}
}

func getCards() []int {
	cards := []int{}

	for i := 1; i <= 13; i++ {
		cards = append(cards, i)
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
			deck.Cards = append(deck.Cards, NewCard(card, suit))
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
