package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type enclosureHandler struct {
	usecase   usecase.EnclosureUsecaseItf
	validator *validator.Validate
}

func RegisterEnclosureHandler(
	usecase usecase.EnclosureUsecaseItf,
	validator *validator.Validate,
	router fiber.Router,
) {
	enclosureHandler := enclosureHandler{usecase, validator}
	router = router.Group("/enclosures")
	router.Get("", enclosureHandler.GetAll)
}

func (h *enclosureHandler) GetAll(c *fiber.Ctx) error {
	enclosures, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(enclosures)
}
