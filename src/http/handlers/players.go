package handlers

import (
	"errors"
	"fmt"
	"qpoker/auth"
	"qpoker/models"

	"github.com/gofiber/fiber"
	"github.com/goware/emailx"
)

type createPlayerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	PW       string `json:"pw"`
}

func (c createPlayerRequest) validate() error {
	if len(c.Username) < 6 {
		return errors.New("username requires at least 6 characters")
	}

	if len(c.PW) < 6 {
		return errors.New("password requires at least 6 characters")
	}

	if err := emailx.Validate(c.Email); err != nil {
		return fmt.Errorf("Invalid email format: %s", c.Email)
	}

	return nil
}

// CreatePlayer creates a new player
func CreatePlayer(c *fiber.Ctx) {
	var req createPlayerRequest

	err := c.BodyParser(&req)
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

	player := &models.Player{
		Username: req.Username,
		Email:    req.Email,
	}

	err = player.Create(req.PW)
	if err != nil {
		c.SendStatus(500)
		c.JSON(formatErrors(err))
		return
	}

	player.Token, err = auth.CreateToken(player)
	if err != nil {
		c.SendStatus(500)
		c.JSON(formatErrors(err))
		return
	}

	c.SendStatus(201)
	c.JSON(player)
}
