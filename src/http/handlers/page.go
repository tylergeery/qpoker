package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"qpoker/models"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/gofiber/fiber"
)

func file(filename string) string {
	return fmt.Sprintf("/src/http/views/%s.mustache", filename)
}

func pageVars(vars fiber.Map) fiber.Map {
	vars["qpoker_host"] = os.Getenv("QPOKER_HOST")

	return vars
}

// RenderPage renders mustache templates
func RenderPage(c *fiber.Ctx, filename string, bind fiber.Map) {
	fp := &mustache.FileProvider{
		Paths:      []string{"", "/src/http/views/"},
		Extensions: []string{"", ".mustache"},
	}

	layoutTmpl, err := mustache.ParseFilePartials(file("layout"), fp)
	if err != nil {
		c.SendStatus(500)
		c.SendString(err.Error())
		return
	}

	tmpl, err := mustache.ParseFilePartials(file(filename), fp)
	if err != nil {
		c.SendStatus(500)
		c.SendString(err.Error())
		return
	}

	html, err := tmpl.RenderInLayout(layoutTmpl, bind)
	if err != nil {
		c.SendStatus(500)
		c.SendString(err.Error())
		return
	}

	c.Set("Content-Type", "text/html")
	c.SendString(html)
}

// PageLanding renders the default app landing page
func PageLanding(c *fiber.Ctx) {
	RenderPage(c, "main", pageVars(fiber.Map{
		"title":       "QCards - Video Social Card Games",
		"stylesheets": []string{"main"},
		"scripts":     []string{"main"},
	}))
}

// PageTable renders a poker table
func PageTable(c *fiber.Ctx) {
	gameSlug := strings.ToLower(c.Params("slug"))

	game, err := models.GetGameBy("slug", gameSlug)
	if err != nil {
		c.SendStatus(404)
		RenderPage(c, "error", pageVars(fiber.Map{}))
		return
	}

	gameObject, err := json.Marshal(game)
	if err != nil {
		c.SendStatus(500)
		RenderPage(c, "error", pageVars(fiber.Map{}))
		return
	}

	RenderPage(c, "table", pageVars(fiber.Map{
		"title":       fmt.Sprintf("%s | QCards Table", game.Name),
		"stylesheets": []string{"game", "games/poker"},
		"scripts":     []string{"games/poker"},
		"game":        string(gameObject),
		"gameOwner":   game.OwnerID,
	}))
}
