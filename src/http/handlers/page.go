package handlers

import (
	"fmt"
	"qpoker/models"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/gofiber/fiber"
)

func file(filename string) string {
	return fmt.Sprintf("/src/http/views/%s.mustache", filename)
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
	RenderPage(c, "main", fiber.Map{
		"title":       "App",
		"stylesheets": []string{"main"},
		"scripts":     []string{},
	})
}

// PageTable renders a poker table
func PageTable(c *fiber.Ctx) {
	gameSlug := strings.ToLower(c.Params("slug"))

	game, err := models.GetGameBy("slug", gameSlug)
	if err != nil {
		c.SendStatus(404)
		RenderPage(c, "error", fiber.Map{})
		return
	}

	RenderPage(c, "table", fiber.Map{
		"title":       fmt.Sprintf("Table %s", game.Name),
		"stylesheets": []string{"table"},
		"scripts":     []string{"table"},
	})
}
