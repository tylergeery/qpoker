package cards

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	SUIT_CLUBS    = "C"
	SUIT_DIAMONDS = "D"
	SUIT_HEARTS   = "H"
	SUIT_SPADES   = "S"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Card is a single card object
type Card struct {
	value int
	suit  string
}

// Deck represent a single deck of cards
type Deck struct {
	Cards []Card
	index int
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
