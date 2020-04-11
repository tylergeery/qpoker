package models

import (
	"fmt"
	"qpoker/utils"
	"time"

	"github.com/metal3d/go-slugify"
)

// Game handles user info
type Game struct {
	ID        int64       `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	OwnerID   int64       `json:"-"`
	Options   GameOptions `json:"options"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetGameBy returns a Game found from the id
func GetGameBy(key string, val interface{}) (*Game, error) {
	game := &Game{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT id, name, slug, owner_id, created_at, updated_at
		FROM game
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(&game.ID, &game.Name, &game.Slug, &game.OwnerID, &game.CreatedAt, &game.UpdatedAt)

	return game, err
}

// Save writes the Game object to the database
func (g *Game) Save() error {
	if g.ID == 0 {
		g.createSlug()
		return g.insert()
	}

	return g.update()
}

func (g *Game) createSlug() {
	g.Slug = fmt.Sprintf("%s-%s", slugify.Marshal(g.Name)[:10], utils.GenerateSlug(5))
}

func (g *Game) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game (name, slug, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`, g.Name, g.Slug, g.OwnerID).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *Game) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game
		SET
			name = $2,
			slug = $3,
			owner_id = $4,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, g.Name, g.Slug, g.OwnerID).Scan(&g.UpdatedAt)

	return err
}
