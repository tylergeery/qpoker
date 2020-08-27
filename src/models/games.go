package models

import (
	"fmt"
	"qpoker/qutils"
	"strings"
	"time"

	"github.com/metal3d/go-slugify"
)

const (
	// GameStatusInit is initial game state
	GameStatusInit = "init"
	// GameStatusStart is started game state
	GameStatusStart = "start"
	// GameStatusComplete is completed game state
	GameStatusComplete = "complete"
)

// Game handles user info
type Game struct {
	ID         int64                  `json:"id"`
	Name       string                 `json:"name"`
	Slug       string                 `json:"slug"`
	OwnerID    int64                  `json:"-"`
	GameTypeID int64                  `json:"game_type_id"`
	Status     string                 `json:"status"`
	Options    map[string]interface{} `json:"options"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// GetGameBy returns a Game found from the id
func GetGameBy(key string, val interface{}) (*Game, error) {
	game := &Game{}

	err := ConnectToDB().QueryRow(fmt.Sprintf(`
		SELECT id, name, slug, owner_id, game_type_id, status, created_at, updated_at
		FROM game
		WHERE %s = $1
		LIMIT 1
	`, key), val).Scan(
		&game.ID, &game.Name, &game.Slug, &game.OwnerID, &game.GameTypeID,
		&game.Status, &game.CreatedAt, &game.UpdatedAt)

	if err != nil {
		return game, err
	}

	gameOptions, _ := GetGameOptionsForGame(game.ID, game.GameTypeID)
	game.Options = gameOptions.Options

	return game, err
}

// IsComplete checks whether the game is complete
func (g *Game) IsComplete() bool {
	return g.Status == GameStatusComplete
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
	slug := slugify.Marshal(g.Name)
	if len(slug) > 10 {
		slug = slug[:10]
	}

	slug = strings.TrimSuffix(slug, "-")
	slug = fmt.Sprintf("%s-%s", slug, qutils.GenerateVariedLengthSlug(5, 15))

	g.Slug = strings.ToLower(slug)
}

func (g *Game) insert() error {
	err := ConnectToDB().QueryRow(`
		INSERT INTO game (name, slug, owner_id, game_type_id, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, g.Name, g.Slug, g.OwnerID, g.GameTypeID, g.Status).Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt)

	return err
}

func (g *Game) update() error {
	err := ConnectToDB().QueryRow(`
		UPDATE game
		SET
			name = $2,
			slug = $3,
			owner_id = $4,
			game_type_id = $5,
			status = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`, g.ID, g.Name, g.Slug, g.OwnerID, g.GameTypeID, g.Status).Scan(&g.UpdatedAt)

	return err
}
