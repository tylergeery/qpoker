package handlers

import (
	"fmt"
	"qpoker/http/utils"
	"qpoker/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber"
)

type createGameRequest struct {
	Name       string                 `json:"name"`
	GameTypeID int64                  `json:"game_type_id"`
	Options    map[string]interface{} `json:"options"`
}

func (req createGameRequest) validate() error {
	if len(req.Name) < 6 {
		return fmt.Errorf("Game name must be at least 6 characters")
	}

	if req.GameTypeID <= int64(0) {
		return fmt.Errorf("Invalid game type id: %d", req.GameTypeID)
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
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown user")))
		return
	}

	err = c.BodyParser(&req)
	if err != nil {
		fmt.Printf("JSON error parsing body (%s): %s\n", c.Fasthttp.Request.Body(), err)
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	err = req.validate()
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	game := &models.Game{
		Name:       req.Name,
		OwnerID:    player.ID,
		Status:     models.GameStatusInit,
		GameTypeID: req.GameTypeID,
	}
	err = game.Save()
	if err != nil {
		fmt.Printf("Game save error: %s", err)
		c.SendStatus(500)
		c.JSON(utils.FormatErrors(err))
		return
	}

	options := &models.GameOptions{
		GameID:     game.ID,
		GameTypeID: game.GameTypeID,
		Options:    req.Options,
	}
	err = options.Save()
	if err != nil {
		fmt.Printf("Game options save error: %s\n", err)
	}

	game, _ = models.GetGameBy("id", game.ID)

	c.SendStatus(201)
	c.JSON(game)
}

type updateGameRequest createGameRequest

func (req updateGameRequest) validate() error {
	if req.Name != "" && len(req.Name) < 6 {
		return fmt.Errorf("Game name must be at least 6 characters")
	}

	return nil
}

// UpdateGame updates an existing game
func UpdateGame(c *fiber.Ctx) {
	var req updateGameRequest

	playerID := c.Locals("playerID").(int64)
	gameID, err := strconv.Atoi(c.Params("gameID"))
	if err != nil {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game ID type")))
		return
	}

	game, err := models.GetGameBy("id", gameID)
	if err != nil || game.OwnerID != playerID {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game")))
		return
	}

	err = c.BodyParser(&req)
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	err = req.validate()
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	if req.Name != "" {
		game.Name = req.Name
	}

	err = game.Save()
	if err != nil {
		fmt.Printf("Game save error: %s", err)
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	options, err := models.GetGameOptionsForGame(game.ID, game.GameTypeID)
	if err != nil {
		options = models.GameOptions{
			GameID:     game.ID,
			GameTypeID: game.GameTypeID,
			Options:    map[string]interface{}{},
		}
	}

	for k, v := range req.Options {
		options.Options[k] = v
	}

	err = options.Save()
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	game, _ = models.GetGameBy("id", game.ID)

	c.SendStatus(200)
	c.JSON(game)
}

// GetGame return a specified game
func GetGame(c *fiber.Ctx) {
	playerID := c.Locals("playerID").(int64)
	gameID, err := strconv.Atoi(c.Params("gameID"))
	if err != nil {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game ID type")))
		return
	}

	game, err := models.GetGameBy("id", gameID)
	if err != nil || game.OwnerID != playerID {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game")))
		return
	}

	c.SendStatus(200)
	c.JSON(game)
}

// GetGameHistory returns all game history events
func GetGameHistory(c *fiber.Ctx) {
	playerID := c.Locals("playerID").(int64)
	gameID, err := strconv.Atoi(c.Params("gameID"))
	if err != nil {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game ID type")))
		return
	}

	game, err := models.GetGameBy("id", gameID)
	if err != nil {
		c.SendStatus(404)
		c.JSON(utils.FormatErrors(fmt.Errorf("Unknown game")))
		return
	}

	since, err := time.Parse(c.Query("since"), time.RFC3339Nano)
	if err != nil {
		since = game.CreatedAt
	}

	if game.IsComplete() {
		// TODO: ensure player actually played
	}

	gameHands, err := models.GetHandsForGame(game.ID, playerID, since, 100)
	if err != nil {
		c.SendStatus(500)
		c.JSON(utils.FormatErrors(fmt.Errorf("Error querying game hands")))
		return
	}

	chipRequests, err := models.GetChipRequestsForGame(game.ID, since, 100)
	if err != nil {
		c.SendStatus(500)
		c.JSON(utils.FormatErrors(fmt.Errorf("Error querying game chip requests")))
		return
	}

	c.SendStatus(200)
	c.JSON(combineAndSort(gameHands, chipRequests))
}

func combineAndSort(hands []*models.GameHandWithPlayer, reqs []*models.GameChipRequest) []interface{} {
	sorted := []interface{}{}

	for len(hands) > 0 || len(reqs) > 0 {
		if len(hands) > 0 && len(reqs) > 0 {
			if hands[0].CreatedAt.Before(reqs[0].UpdatedAt) {
				sorted = append(sorted, hands[0])
				hands = hands[1:]
			} else {
				sorted = append(sorted, reqs[0])
				reqs = reqs[1:]
			}

			continue
		}

		if len(hands) > 0 {
			sorted = append(sorted, hands[0])
			hands = hands[1:]
			continue
		}

		sorted = append(sorted, reqs[0])
		reqs = reqs[1:]
	}

	return sorted
}
