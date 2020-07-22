package hearts

import (
	"qpoker/cards"
	"qpoker/models"
)

// Player holds the information about a player at a table
type Player struct {
	ID       int64           `json:"id"`
	Username string          `json:"username"`
	Cards    []cards.Card    `json:"-"`
	Pile     []cards.Card    `json:"pile"`
	Options  map[string]bool `json:"options"`
	Score    int64
}

// NewPlayer creates a new Player
func NewPlayer(player *models.Player) *Player {
	return &Player{
		ID:       player.ID,
		Username: player.Username,
	}
}

// GetID return player ID
func (p *Player) GetID() int64 {
	return p.ID
}

// SetPlayerActions sets the moves a player is allowed to make
func (p *Player) SetPlayerActions(actions map[string]bool) {
	p.Options = actions
}

// IsActive returns whether the player is active in the current hand
func (p *Player) IsActive() bool {
	return true
}

// IsReady returns whether the player is ready for next game
func (p *Player) IsReady() bool {
	return true
}

// AddCards adds cards to players hand
func (p *Player) AddCards(c []cards.Card) {
	p.Cards = append(p.Cards, c...)
	cards.SortSuitedAcesHigh(p.Cards)
}

// RemoveCards removes cards from players hand
func (p *Player) RemoveCards(c []cards.Card) {
	cardMap := map[string]bool{}
	for _, card := range c {
		cardMap[card.ToString()] = true
	}

	for i := 0; ; {
		if i >= len(p.Cards) {
			break
		}

		card := p.Cards[i]
		if _, ok := cardMap[card.ToString()]; !ok {
			i++
			continue
		}

		p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)
	}
}

// HeartsCount gets hearts count in pile
func (p *Player) HeartsCount() int64 {
	count := int64(0)
	for _, c := range p.Pile {
		if c.Suit == cards.SuitHearts {
			count++
		}
		if c.Suit == cards.SuitSpades && c.Value == 12 {
			count += 13
		}
	}

	return count
}
