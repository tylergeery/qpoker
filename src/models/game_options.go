package models

import (
	"fmt"
	"time"

	"database/sql/driver"
	"encoding/json"
)

// GameOptions is a config object used to control game settings
type GameOptions struct {
	Capacity         int   `json:"capacity"`
	BigBlind         int64 `json:"big_blind"`
	TimeBetweenHands int   `json:"time_between_hands"`
	BuyInMin         int64 `json:"buy_in_min"`
	BuyInMax         int64 `json:"buy_in_max"`
}

// Value makes the GameOptions struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (g GameOptions) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// Scan makes the GameOptions struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (g *GameOptions) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Could not get GameOptions as []byte")
	}

	return json.Unmarshal(b, &g)
}

// GameOptionsRecord handles user info
type GameOptionsRecord struct {
	ID        int64       `json:"id"`
	GameID    int64       `json:"game_id"`
	Options   GameOptions `json:"options"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetGameOptionsForGame returns options for game
func GetGameOptionsForGame(gameID int64) (GameOptions, error) {
	record, err := GetGameOptionsRecordBy("game_id", gameID)
	if err != nil {
		return GameOptions{}, err
	}

	return record.Options, nil
}

// GetGameOptionsRecordBy returns a GameOptionsRecord found from key,val
func GetGameOptionsRecordBy(key string, val interface{}) (*GameOptionsRecord, error) {
	record := &GameOptionsRecord{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT id, game_id, options, created_at, updated_at
		FROM game_options
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(&record.ID, &record.GameID, &record.Options, &record.CreatedAt, &record.UpdatedAt)

	if record.ID == 0 {
		return record, fmt.Errorf("GameOptionsRecord could not be found for %s=%s", key, val)
	}

	return record, err
}

// Save writes the Game object to the database
func (g *GameOptionsRecord) Save() error {
	if g.ID == 0 {
		return g.insert()
	}

	return g.update()
}

func (g *GameOptionsRecord) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game_options (game_id, options)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`, g.GameID, g.Options).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *GameOptionsRecord) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game_options
		SET options = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, g.Options).Scan(&g.UpdatedAt)

	return err
}
