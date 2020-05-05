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
	Name    string             `json:"name"`
	Options models.GameOptions `json:"options"`
}

func (req createGameRequest) validate() error {
	if len(req.Name) < 6 {
		return fmt.Errorf("Game name must be at least 6 characters")
	}

	// Allow -1 as an infinite capacity
	if req.Options.Capacity == 0 || req.Options.Capacity < -1 {
		return fmt.Errorf("Invalid capacity: %d", req.Options.Capacity)
	}

	if req.Options.BigBlind <= 0 || req.Options.BigBlind%2 == 1 {
		return fmt.Errorf("Invalid big blind: %d", req.Options.BigBlind)
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
		Name:    req.Name,
		OwnerID: player.ID,
		Status:  models.GameStatusInit,
	}
	err = game.Save()
	if err != nil {
		fmt.Printf("Game save error: %s", err)
		c.SendStatus(500)
		c.JSON(utils.FormatErrors(err))
		return
	}

	// Sane defaults
	req.Options.TimeBetweenHands = 5
	req.Options.BuyInMin = req.Options.BigBlind
	req.Options.BuyInMax = req.Options.BigBlind * 10

	options := &models.GameOptionsRecord{
		GameID:  game.ID,
		Options: req.Options,
	}
	err = options.Save()
	if err != nil {
		fmt.Printf("Game options save error: %s\n", err)
	}
	fmt.Printf("gameOptions after save: %+v\n", options)
	if err == nil {
		game.Options = options.Options
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
	if req.Options.Capacity < -1 {
		return fmt.Errorf("Invalid capacity: %d", req.Options.Capacity)
	}

	if req.Options.BigBlind < 0 || req.Options.BigBlind%2 == 1 {
		return fmt.Errorf("Invalid big blind: %d", req.Options.BigBlind)
	}

	if req.Options.TimeBetweenHands < 0 || req.Options.TimeBetweenHands > 30 {
		return fmt.Errorf("Invalid time between hands: %d", req.Options.TimeBetweenHands)
	}

	if req.Options.BuyInMin < 0 {
		return fmt.Errorf("Invalid min buy in: %d", req.Options.BuyInMin)
	}

	if req.Options.BuyInMax < 0 {
		return fmt.Errorf("Invalid max buy in: %d", req.Options.BuyInMax)
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

	if req.Options.Capacity != 0 {
		game.Options.Capacity = req.Options.Capacity
	}
	if req.Options.BigBlind != 0 {
		game.Options.BigBlind = req.Options.BigBlind
	}
	if req.Options.TimeBetweenHands != 0 {
		game.Options.TimeBetweenHands = req.Options.TimeBetweenHands
	}
	if req.Options.BuyInMin != 0 {
		game.Options.BuyInMin = req.Options.BuyInMin
	}
	if req.Options.BuyInMax != 0 {
		game.Options.BuyInMax = req.Options.BuyInMax
	}

	if req.Options.BuyInMax < req.Options.BuyInMin {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(fmt.Errorf("Game buy in max cannot be less than min")))
		return
	}

	if req.Name != "" {
		game.Name = req.Name
	}

	err = game.Save()
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

	optionsRecord, err := models.GetGameOptionsRecordBy("game_id", game.ID)
	if err != nil {
		optionsRecord = &models.GameOptionsRecord{GameID: game.ID}
	}

	optionsRecord.Options = game.Options
	err = optionsRecord.Save()
	if err != nil {
		c.SendStatus(400)
		c.JSON(utils.FormatErrors(err))
		return
	}

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
		} else {
			sorted = append(sorted, reqs[0])
			reqs = reqs[1:]
		}
	}

	return sorted
}
