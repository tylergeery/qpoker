package handlers

import (
	"errors"
	"fmt"
	"qpoker/auth"
	"qpoker/models"

	"github.com/gofiber/fiber"
)

type createPlayerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	pw       string `json:"pw"`
}

func (c createPlayerRequest) validate() error {
	if len(c.Username) < 6 {
		return errors.New("username requires at least 6 characters")
	}

	if len(c.pw) < 6 {
		return errors.New("password requires at least 6 characters")
	}

	// TODO: validate email

	return nil
}

// CreatePlayer creates a new player
func CreatePlayer(c *fiber.Ctx) {
	var req createPlayerRequest

	err := c.BodyParser(&req)
	if err != nil {
		c.SendStatus(400)
		c.Send(fmt.Sprintf("Bad Request: %s", err))
		return
	}

	err = req.validate()
	if err != nil {
		c.SendStatus(400)
		c.Send(fmt.Sprintf("Bad Request: %s", err))
		return
	}

	player := &models.Player{
		Username: req.Username,
		Email:    req.Email,
	}

	err = player.Create(req.pw)
	if err != nil {
		c.SendStatus(500)
		c.Send(fmt.Sprintf("Internal Error: %s", err))
		return
	}

	player.Token, err = auth.CreateToken(player)
	if err != nil {
		c.SendStatus(500)
		c.Send(fmt.Sprintf("Internal Error: %s", err))
		return
	}

	c.JSON(player)
}
