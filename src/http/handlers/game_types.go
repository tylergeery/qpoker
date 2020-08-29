package handlers

import (
	"fmt"
	"qpoker/http/utils"
	"qpoker/models"

	"github.com/gofiber/fiber"
)

// GameType contains game types and otions for api response
type GameType struct {
	ID          int64               `json:"id"`
	Key         string              `json:"key"`
	DisplayName string              `json:"display_name"`
	Options     []models.GameOption `json:"options"`
}

// GetGameTypes gets all game types and options
func GetGameTypes(c *fiber.Ctx) {
	response := []GameType{}
	gameTypes, err := models.GetGameTypes()
	if err != nil {
		c.SendStatus(500)
		c.JSON(utils.FormatErrors(err))
		return
	}

	for _, gameType := range gameTypes {
		options, err := models.GetGameOptionRecordsForGameType(gameType.ID)
		if err != nil {
			fmt.Printf("Game Type options fetch err: %s\n", err)
			continue
		}

		gt := GameType{gameType.ID, gameType.Key, gameType.DisplayName, options}
		response = append(response, gt)
	}

	c.SendStatus(200)
	c.JSON(response)
}
