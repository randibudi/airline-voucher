package app

import "github.com/gofiber/fiber/v3"

// New creates and configures the HTTP application.
func New() *fiber.App {
	application := fiber.New()

	application.Get("/api/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return application
}
