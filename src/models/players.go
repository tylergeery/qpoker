package models

import (
	"time"
)

// Player handles user info
type Player struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	pw        string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetPlayerFromID returns a Player found from the id
func GetPlayerFromID(id int64) (*Player, error) {
	player := &Player{}

	err := ConnectToDB().QueryRow(`
		SELECT id, username, email, created_at, updated_at
		FROM player
		WHERE id = $1
		LIMIT 1
	`, id).Scan(&player.ID, &player.Username, &player.Email, &player.CreatedAt, &player.UpdatedAt)

	return player, err
}

// Save the last player info
func (p *Player) Save() error {
	return p.update()
}

// Create a new player
func (p *Player) Create(hashedPW string) error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO player (username, email, pw)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, p.Username, p.Email, hashedPW).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	return err
}

func (p *Player) update() error {
	rows, err := ConnectToDB().Query(`
		UPDATE player
		SET username = $2, email = $3, updated_at = NOW()
		WHERE id = $1
	`, p.ID, p.Username, p.Email)
	defer rows.Close()

	return err
}
