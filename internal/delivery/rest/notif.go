package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type notifHandler struct {
	usecase usecase.NotifUsecaseItf
}

func RegisterNotifHandler(
	usecase usecase.NotifUsecaseItf,
	router fiber.Router,
) {
	notifHandler := notifHandler{usecase}
	router = router.Group("/notifications")
	router.Get("", notifHandler.notifications)
}

func (h *notifHandler) notifications(c *fiber.Ctx) error {
	notifs, err := h.usecase.Fetch(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"notifications": notifs})
}
