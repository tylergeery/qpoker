package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// UserStackMap hold all maps of userID to int64 stack
type UserStackMap map[int64]int64

// Value makes the UserStackMap struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (u UserStackMap) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}

	return json.Marshal(u)
}

// Scan makes the UserStackMap struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (u *UserStackMap) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Could not get UserStackMap as []byte")
	}

	return json.Unmarshal(b, &u)
}

// GameHand holds a single game hand
type GameHand struct {
	ID        int64        `json:"id"`
	GameID    int64        `json:"game_id"`
	Board     []string     `json:"board"`
	Payouts   UserStackMap `json:"payouts"`
	Bets      UserStackMap `json:"bets"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// GameHandWithPlayer returns game hand with player info
type GameHandWithPlayer struct {
	GameHand
	Cards    []string      `json:"cards"`
	Starting sql.NullInt64 `json:"starting"`
	Ending   sql.NullInt64 `json:"ending"`
}

// GetGameHandBy returns a GameHand found from key,val
func GetGameHandBy(key string, val interface{}) (*GameHand, error) {
	hand := &GameHand{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT
			id,
			game_id,
			board,
			payouts,
			bets,
			created_at,
			updated_at
		FROM game_hand
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(
		&hand.ID, &hand.GameID, pq.Array(&hand.Board), &hand.Payouts,
		&hand.Bets, &hand.CreatedAt, &hand.UpdatedAt,
	)

	return hand, err
}

// GetHandsForGame returns all hands for game
func GetHandsForGame(gameID, playerID int64, since time.Time, count int) ([]*GameHandWithPlayer, error) {
	hands := []*GameHandWithPlayer{}

	rows, err := ConnectToDB().Query(`
		SELECT
			gh.id,
			gh.game_id,
			board,
			gh.payouts,
			gh.bets,
			gh.created_at,
			gh.updated_at,
			gph.cards,
			gph.starting,
			gph.ending
		FROM game_hand gh
		LEFT JOIN game_player_hand gph ON (
			gh.id = gph.game_hand_id AND gph.player_id = $3
		)
		WHERE game_id = $1
			AND gh.created_at > $2
		ORDER BY gh.created_at ASC
		LIMIT $4
	`, gameID, since, playerID, count)
	if err != nil {
		return hands, err
	}

	defer rows.Close()
	for rows.Next() {
		hand := &GameHandWithPlayer{}

		rows.Scan(
			&hand.ID, &hand.GameID, pq.Array(&hand.Board), &hand.Payouts,
			&hand.Bets, &hand.CreatedAt, &hand.UpdatedAt,
			pq.Array(&hand.Cards), &hand.Starting, &hand.Ending,
		)
		hands = append(hands, hand)
	}

	return hands, nil
}

// Save writes the Game object to the database
func (g *GameHand) Save() error {
	if g.ID == 0 {
		return g.insert()
	}

	return g.update()
}

func (g *GameHand) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game_hand (game_id)
		VALUES ($1)
		RETURNING id, created_at, updated_at
	`, g.GameID).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *GameHand) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game_hand
		SET
			board = $2,
			payouts = $3,
			bets = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, pq.StringArray(g.Board), g.Payouts, g.Bets).Scan(&g.UpdatedAt)

	return err
}
