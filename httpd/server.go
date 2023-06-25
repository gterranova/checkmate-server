package httpd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"terra9.it/vadovia/assets"
	"terra9.it/vadovia/core"
)

func ProjectData(ctx *fiber.Ctx) error {
	project := core.NewProject()
	return ctx.JSON(
		fiber.Map{
			"project": project,
			"result":  project.Describe(),
		},
	)
}

func SetupRoutes(app *fiber.App) {

	app.Get("/assets/logo.png", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/png")
		return c.Send(assets.ResourceLogovadoviacompactPng.StaticContent)
	})

	route := app.Group("/api/v1")

	route.Get("/project", ProjectData)
}

func New(serverAddr string) {

	app := fiber.New()

	// Add CORS Middleware so the frontend get the cookie
	app.Use(cors.New())

	SetupRoutes(app)

	app.Listen(serverAddr)
}
