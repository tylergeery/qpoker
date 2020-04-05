package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	return getPlayerByKey("id", id)
}

// GetPlayerFromEmail returns a Player found from the id
func GetPlayerFromEmail(email string) (*Player, error) {
	return getPlayerByKey("email", email)
}

// AuthenticatePlayer returns a player if authenticated, otherwise an error
func AuthenticatePlayer(email, pw string) (*Player, error) {
	player, err := GetPlayerFromEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(player.pw), []byte(pw))
	if err != nil {
		return nil, err
	}

	return player, nil
}

func getPlayerByKey(key string, val interface{}) (*Player, error) {
	player := &Player{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT id, username, email, pw, created_at, updated_at
		FROM player
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(&player.ID, &player.Username, &player.Email, &player.pw, &player.CreatedAt, &player.UpdatedAt)

	return player, err
}

// Save the last player info
func (p *Player) Save() error {
	return p.update()
}

// Create a new player
func (p *Player) Create(pw string) error {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = ConnectToDB().QueryRow(`
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
