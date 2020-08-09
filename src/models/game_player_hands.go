package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// GamePlayerHand holds a single game hand
type GamePlayerHand struct {
	ID         int64     `json:"id"`
	GameHandID int64     `json:"game_hand_id"`
	PlayerID   int64     `json:"player_id"`
	Cards      []string  `json:"cards"`
	Starting   int64     `json:"starting"`
	Ending     int64     `json:"ending"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetGamePlayerHandBy returns a GamePlayerHand found from key,val
func GetGamePlayerHandBy(key string, val interface{}) (*GamePlayerHand, error) {
	var ending sql.NullInt64
	playerHand := &GamePlayerHand{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT
			id, game_hand_id, player_id,
			cards, starting, ending,
			created_at, updated_at
		FROM game_player_hand
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(
		&playerHand.ID, &playerHand.GameHandID, &playerHand.PlayerID, pq.Array(&playerHand.Cards),
		&playerHand.Starting, &ending,
		&playerHand.CreatedAt, &playerHand.UpdatedAt,
	)

	if ending.Valid {
		playerHand.Ending = ending.Int64
	}

	return playerHand, err
}

// GetGamePlayerHandForGameAndPlayer returns a GamePlayerHand found for game_id
func GetGamePlayerHandForGameAndPlayer(gameID, playerID int64) (*GamePlayerHand, error) {
	var ending sql.NullInt64
	playerHand := &GamePlayerHand{}

	err := ConnectToDB().QueryRow(`
		SELECT
			gph.id, gph.game_hand_id, gph.player_id,
			gph.cards, gph.starting, gph.ending,
			gph.created_at, gph.updated_at
		FROM game_player_hand gph
		JOIN game_hand gh ON (gh.id = gph.game_hand_id)
		WHERE gh.game_id = $1
			AND gph.player_id = $2
		ORDER BY gh.created_at DESC
		LIMIT 1
	`, gameID, playerID).Scan(
		&playerHand.ID, &playerHand.GameHandID, &playerHand.PlayerID, pq.Array(&playerHand.Cards),
		&playerHand.Starting, &ending,
		&playerHand.CreatedAt, &playerHand.UpdatedAt,
	)

	playerHand.Ending = int64(-1)
	if ending.Valid {
		playerHand.Ending = ending.Int64
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
		INSERT INTO game_player_hand (game_hand_id, player_id, cards, starting, ending)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, g.GameHandID, g.PlayerID, pq.StringArray(g.Cards), g.Starting, g.Ending).Scan(
		&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *GamePlayerHand) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game_player_hand
		SET
			cards = $2,
			ending = $3,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, pq.StringArray(g.Cards), g.Ending).Scan(&g.UpdatedAt)

	return err
}
