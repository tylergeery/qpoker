package hearts

import (
	"fmt"
	"qpoker/cards"
)

const (
	// PassLeft means pass left
	PassLeft = 'l'
	// PassRight means pass right
	PassRight = 'r'
	// PassAcross means pass across
	PassAcross = 'a'
	// PassNone means no pass
	PassNone = 'n'

	// StatePassing game state when all players are selecting cards to pass
	StatePassing = "Passing"
	// StateActive game state after game is active
	StateActive = "Active"
)

var passDirections []byte = []byte{PassLeft, PassRight, PassAcross, PassNone}

// Hearts contains game state
type Hearts struct {
	Board  map[int64]cards.Card `json:"board"`
	Piles  map[int64]cards.Card `json:"piles"`
	Deck   *cards.Deck          `json:"-"`
	Table  *Table               `json:"table"`
	State  string               `json:"state"`
	Scores map[int64]int        `json:"scores"`

	passes    map[int64][]cards.Card
	passIndex int
}

// NewHearts creates and returns the resources for a new Hearts game
func NewHearts(table *Table, state string) *Hearts {
	scores := map[int64]int{}
	for _, player := range table.GetAllPlayers() {
		scores[player.ID] = 0
	}

	return &Hearts{
		Deck:      cards.NewDeck(),
		Table:     table,
		State:     state,
		Scores:    scores,
		passIndex: -1,
	}
}

// Deal a new hand of hearts
func (h *Hearts) Deal() error {
	// reset table
	h.Table.NextHand()

	// reset board
	h.State = StatePassing
	h.Board = map[int64]cards.Card{}
	h.passes = map[int64][]cards.Card{}
	h.passIndex = (h.passIndex + 1) % len(passDirections)

	// shuffle cards
	h.Deck.Shuffle()

	// deal player cards
	next := h.Table.ActiveIndex
	players := h.Table.GetActivePlayers()

	for i := 0; i < 52; i++ {
		card, err := h.Deck.GetCard()
		if err != nil {
			return fmt.Errorf("unexpected GetCard deal error hearts: %s", err)
		}

		players[next].Cards = append(players[next].Cards, card)
		next = h.Table.NextPos(next)
	}

	if h.skipPass() {
		h.State = StateActive
	}

	return nil
}

func (h *Hearts) clearBoard() {
	h.Board = map[int64]cards.Card{}
}

func (h *Hearts) playerPlay(playerID int64, card cards.Card) {
	h.Board[playerID] = card

	// Remove card from players
	player := h.Table.GetPlayerByID(playerID)
	for i := range player.Cards {
		if player.Cards[i].ToString() == card.ToString() {
			player.Cards = append(player.Cards[:i], player.Cards[i+1:]...)
			break
		}
	}
}

func (h *Hearts) addPass(playerID int64, cards []cards.Card) {
	h.passes[playerID] = cards

	if len(h.passes) == len(h.Table.Players) {
		// commit passes
		h.pass()
		h.State = StateActive
	}
}

func (h *Hearts) getPassToPlayer(playerPos int) *Player {
	var passPlayerPos int

	switch passDirections[h.passIndex] {
	case PassLeft:
		passPlayerPos = h.Table.NextPos(h.Table.NextPos(h.Table.NextPos(playerPos)))
	case PassRight:
		passPlayerPos = h.Table.NextPos(playerPos)
	case PassAcross:
		passPlayerPos = h.Table.NextPos(h.Table.NextPos(playerPos))
	}

	return h.Table.GetPlayerByID(h.Table.GetPlayerID(passPlayerPos))
}

func (h *Hearts) pass() {
	if h.State != StatePassing {
		return
	}

	start := h.Table.ActiveIndex
	curr := start
	for {
		player, passPlayer := h.Table.GetPlayerByID(h.Table.GetPlayerID(curr)), h.getPassToPlayer(curr)
		player.RemoveCards(h.passes[player.ID])
		passPlayer.AddCards(h.passes[player.ID])

		curr = h.Table.NextPos(curr)
		if curr == start {
			break
		}
	}
}

func (h *Hearts) skipPass() bool {
	return passDirections[h.passIndex] == PassNone
}

func (h *Hearts) passesComplete() bool {
	return h.State == StateActive
}

// PointTotals calculates player point totals
func (h *Hearts) PointTotals(players []*Player) map[int64]int64 {
	totals := map[int64]int64{}
	shooter := int64(0)

	for i := range players {
		playerHearts := players[i].HeartsCount()
		totals[players[i].ID] = playerHearts
		if playerHearts == 26 {
			shooter = players[i].ID
		}
	}

	if shooter == int64(0) {
		return totals
	}

	for i := range players {
		totals[players[i].ID] = 26
		if players[i].ID == shooter {
			totals[players[i].ID] = 0
		}
	}

	return totals
}

// CleanPile collects cards and adds them to player's pile
func (h *Hearts) CleanPile(resetBoard bool) error {
	if len(h.Board) != 4 {
		return fmt.Errorf("Cannot clean pile, board doesn't have expected cards: %+v", h.Board)
	}

	players := h.Table.GetActivePlayers()
	suit, value := h.Board[players[0].ID].Suit, h.Board[players[0].ID].Value
	winner := players[0].ID
	for i := range players {
		playedCard := h.Board[players[i].ID]
		if playedCard.Suit == suit && playedCard.Value > value {
			value, winner = playedCard.Value, players[i].ID
		}
	}

	winningPlayer := h.Table.GetPlayerByID(winner)
	for _, card := range h.Board {
		winningPlayer.Pile = append(winningPlayer.Pile, card)
	}

	if resetBoard {
		h.Board = map[int64]cards.Card{}
	}

	return nil
}
