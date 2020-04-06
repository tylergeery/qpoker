package handlers

import (
	"fmt"
	"qpoker/models"

	"github.com/gofiber/fiber"
)

type createGameRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

func (req createGameRequest) validate() error {
	if len(req.Name) < 6 {
		return fmt.Errorf("Game name must be at least 6 characters")
	}

	// Allow -1 as an infinite capacity
	if req.Capacity == 0 || req.Capacity < -1 {
		return fmt.Errorf("Invalid capacity: %d", req.Capacity)
	}

	return nil
}

// CreateGame creates a new game
func CreateGame(c *fiber.Ctx) {
	var req createGameRequest

	playerID := c.Locals("playerID").(int64)
	player, err := models.GetPlayerFromID(playerID)
	if err != nil {
		c.SendStatus(403)
		c.JSON(formatErrors(fmt.Errorf("Unknown user")))
		return
	}

	err = c.BodyParser(&req)
	if err != nil {
		c.SendStatus(400)
		c.JSON(formatErrors(err))
		return
	}

	err = req.validate()
	if err != nil {
		c.SendStatus(400)
		c.JSON(formatErrors(err))
		return
	}

	game := &models.Game{
		Name:     req.Name,
		OwnerID:  player.ID,
		Capacity: req.Capacity,
	}
	err = game.Save()
	if err != nil {
		c.SendStatus(500)
		c.JSON(formatErrors(err))
		return
	}

	c.SendStatus(201)
	c.JSON(game)
}

type updateGameRequest createGameRequest

func (req updateGameRequest) validate() error {
	if req.Name != "" && len(req.Name) < 6 {
		return fmt.Errorf("Game name must be at least 6 characters")
	}

	// Allow -1 as an infinite capacity
	if req.Capacity < -1 {
		return fmt.Errorf("Invalid capacity: %d", req.Capacity)
	}

	return nil
}

// UpdateGame updates an existing game
func UpdateGame(c *fiber.Ctx) {

}

// GetGame return a specified game
func GetGame(c *fiber.Ctx) {

}
