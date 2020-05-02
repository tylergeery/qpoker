package models

import (
	"fmt"
	"time"
)

const (
	// GameChipRequestStatusInit = "init"
	GameChipRequestStatusInit = "init"
	// GameChipRequestStatusIApproved = "approved"
	GameChipRequestStatusIApproved = "approved"
	// GameChipRequestStatusDenied = "denied"
	GameChipRequestStatusDenied = "denied"
)

// GameChipRequest handles user info
type GameChipRequest struct {
	ID        int64     `json:"id"`
	GameID    int64     `json:"game_id"`
	PlayerID  int64     `json:"player_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetGameChipRequestBy returns a GameChipRequest found from the key
func GetGameChipRequestBy(key string, val interface{}) (*GameChipRequest, error) {
	req := &GameChipRequest{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT id, game_id, player_id, amount, status, created_at, updated_at
		FROM game_chip_requests
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(
		&req.ID, &req.GameID, &req.PlayerID, &req.Amount,
		&req.Status, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return req, err
	}

	return req, err
}

// Save writes the Game object to the database
func (req *GameChipRequest) Save() error {
	if req.ID == 0 {
		return req.insert()
	}

	return req.update()
}

func (req *GameChipRequest) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game_chip_requests (game_id, player_id, amount, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, req.GameID, req.PlayerID, req.Amount, req.Status).Scan(
		&req.ID, &req.CreatedAt, &req.UpdatedAt)

	return err
}

func (req *GameChipRequest) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game
		SET
			amount = $2,
			status = $3,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, req.ID, req.Amount, req.Status).Scan(&req.UpdatedAt)

	return err
}
