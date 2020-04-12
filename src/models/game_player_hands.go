package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// GamePlayerHand holds a single game hand
type GamePlayerHand struct {
	ID            int64     `json:"id"`
	GameHandID    int64     `json:"game_hand_id"`
	PlayerID      int64     `json:"player_id"`
	Cards         []string  `json:"cards"`
	CardsVisible  bool      `json:"cards_visible"`
	StartingStack int64     `json:"starting_stack"`
	EndingStack   int64     `json:"ending_stack"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetGamePlayerHandBy returns a GameHand found from key,val
func GetGamePlayerHandBy(key string, val interface{}) (*GamePlayerHand, error) {
	var endingStack sql.NullInt64
	playerHand := &GamePlayerHand{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT
			id, game_hand_id, player_id, cards,
			starting_stack, ending_stack,
			created_at, updated_at
		FROM game_player_hand
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(
		&playerHand.ID, &playerHand.GameHandID, &playerHand.PlayerID, pq.Array(&playerHand.Cards),
		&playerHand.StartingStack, &endingStack,
		&playerHand.CreatedAt, &playerHand.UpdatedAt,
	)

	if endingStack.Valid {
		playerHand.EndingStack = endingStack.Int64
	}

	return playerHand, err
}

// Save writes the GamePlayerHand object to the database
func (g *GamePlayerHand) Save() error {
	if g.ID == 0 {
		return g.insert()
	}

	return g.update()
}

func (g *GamePlayerHand) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game_player_hand (game_hand_id, player_id, cards, starting_stack)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, g.GameHandID, g.PlayerID, pq.StringArray(g.Cards), g.StartingStack).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *GamePlayerHand) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game_player_hand
		SET
			cards = $2,
			ending_stack = $3,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, pq.StringArray(g.Cards), g.EndingStack).Scan(&g.UpdatedAt)

	return err
}
