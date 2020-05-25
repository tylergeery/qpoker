package models

import (
	"time"

	"qpoker/models/sql"
)

const (
	// GameChipRequestStatusInit = "init"
	GameChipRequestStatusInit = "init"
	// GameChipRequestStatusApproved = "approved"
	GameChipRequestStatusApproved = "approved"
	// GameChipRequestStatusDenied = "denied"
	GameChipRequestStatusDenied = "denied"

	tableName = "game_chip_requests"
)

var (
	columns = []string{
		"id", "game_id", "player_id", "amount", "status",
		"created_at", "updated_at",
	}
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

	query := sql.NewSelect(
		tableName,
		columns,
	).Filter(
		sql.Equal(key, val),
	).Limit(1)

	sql, values := query.ToSQL()
	row := ConnectToDB().QueryRow(sql, values...)
	err := row.Scan(&req.ID, &req.GameID, &req.PlayerID, &req.Amount,
		&req.Status, &req.CreatedAt, &req.UpdatedAt)

	if err != nil {
		return req, err
	}

	return req, err
}

// GetChipRequestsForGame returns all hands for game
func GetChipRequestsForGame(gameID int64, since time.Time, count int) ([]*GameChipRequest, error) {
	query := sql.NewSelect(
		tableName,
		columns,
	).Filter(
		sql.And(
			sql.Equal("game_id", gameID),
			sql.GTE("updated_at", since),
		),
	).Limit(count)

	return getChipRequestsFor(query)
}

// GetApprovedChipRequestsForGameAndPlayer returns all hands for game
func GetApprovedChipRequestsForGameAndPlayer(gameID, playerID int64, since time.Time, count int) ([]*GameChipRequest, error) {
	query := sql.NewSelect(tableName, columns)
	query = query.Filter(
		sql.And(
			sql.Equal("game_id", gameID),
			sql.Equal("player_id", playerID),
			sql.GTE("updated_at", since),
		),
	)
	query = query.Limit(count)

	return getChipRequestsFor(query)
}

func getChipRequestsFor(query sql.Select) ([]*GameChipRequest, error) {
	requests := []*GameChipRequest{}
	sql, values := query.ToSQL()
	rows, err := ConnectToDB().Query(sql, values...)

	if err != nil {
		return requests, err
	}

	defer rows.Close()
	for rows.Next() {
		req := &GameChipRequest{}

		rows.Scan(
			&req.ID, &req.GameID, &req.PlayerID, &req.Amount,
			&req.Status, &req.CreatedAt, &req.UpdatedAt)
		requests = append(requests, req)
	}

	return requests, nil
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
		VALUES ($1, $2, $3, $4)
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
