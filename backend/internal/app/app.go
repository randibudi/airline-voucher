package app

import (
	"github.com/gofiber/fiber/v3"
	"github.com/randibudi/airline-voucher/backend/internal/httpapi"
)

// New creates and configures the HTTP application.
func New(handler *httpapi.Handler) *fiber.App {
	application := fiber.New()

	application.Get("/api/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
	if handler != nil {
		application.Post("/api/check", handler.Check)
		application.Post("/api/generate", handler.Generate)
	}

	return application
}
