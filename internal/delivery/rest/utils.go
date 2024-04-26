package rest

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/pkg/conf"
)

func RegisterUtilsHandler(app fiber.Router) {
	app.Get("/health", healthCheck)
	app.Get("/status", healthCheck)
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "Healthy",
		"uptime": time.Since(conf.StartTime).String(),
	})
}
