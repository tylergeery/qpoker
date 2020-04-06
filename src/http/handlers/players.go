package handlers

import (
	"errors"
	"fmt"
	"qpoker/auth"
	"qpoker/models"
	"strconv"

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

type updatePlayerRequest createPlayerRequest

func (req updatePlayerRequest) validate() error {
	if req.Username != "" && len(req.Username) < 6 {
		return errors.New("username requires at least 6 characters")
	}

	if req.PW != "" && len(req.PW) < 6 {
		return errors.New("password requires at least 6 characters")
	}

	err := emailx.Validate(req.Email)
	if req.Email != "" && err != nil {
		return fmt.Errorf("Invalid email format: %s", req.Email)
	}

	return nil
}

// UpdatePlayer updates player information
func UpdatePlayer(c *fiber.Ctx) {
	var req updatePlayerRequest

	id, err := strconv.Atoi(c.Params("id"))
	playerID := int64(id)

	if err != nil {
		c.SendStatus(404)
		c.JSON(formatErrors(fmt.Errorf("Uknown player ID type")))
		return
	}

	if c.Locals("playerID") != playerID {
		c.SendStatus(403)
		c.JSON(formatErrors(fmt.Errorf("Not authorized to view user")))
		return
	}

	player, err := models.GetPlayerFromID(playerID)
	if err != nil {
		c.SendStatus(403)
		c.JSON(formatErrors(fmt.Errorf("Not authorized to view user")))
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

	if req.Username != "" {
		player.Username = req.Username
	}
	if req.Email != "" {
		player.Email = req.Email
	}
	if req.PW != "" {
		player.SetPassword(req.PW)
	}

	err = player.Save()
	if err != nil {
		c.SendStatus(500)
		c.JSON(formatErrors(err))
		return
	}

	player.Token = c.Locals("token").(string)

	c.SendStatus(200)
	c.JSON(player)
}

type playerLoginRequest struct {
	Email string `json:"email"`
	PW    string `json:"pw"`
}

// LoginPlayer creates a new player
func LoginPlayer(c *fiber.Ctx) {
	var req playerLoginRequest

	err := c.BodyParser(&req)
	if err != nil {
		c.SendStatus(400)
		c.JSON(formatErrors(err))
		return
	}

	player, err := models.AuthenticatePlayer(req.Email, req.PW)
	if err != nil {
		c.SendStatus(400)
		c.JSON(formatErrors(err))
		return
	}

	player.Token, err = auth.CreateToken(player)
	if err != nil {
		c.SendStatus(500)
		c.JSON(formatErrors(err))
		return
	}

	c.SendStatus(200)
	c.JSON(player)
}

// GetPlayer returns the player object
func GetPlayer(c *fiber.Ctx) {
	id, err := strconv.Atoi(c.Params("id"))
	playerID := int64(id)

	if err != nil {
		c.SendStatus(404)
		c.JSON(formatErrors(fmt.Errorf("Uknown player ID type")))
		return
	}

	if c.Locals("playerID") != playerID {
		c.SendStatus(403)
		c.JSON(formatErrors(fmt.Errorf("Not authorized to view user")))
		return
	}

	player, err := models.GetPlayerFromID(playerID)
	if err != nil {
		c.SendStatus(403)
		c.JSON(formatErrors(fmt.Errorf("Not authorized to view user")))
		return
	}

	player.Token = c.Locals("token").(string)

	c.SendStatus(200)
	c.JSON(player)
}
