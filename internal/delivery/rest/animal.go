package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type animalHandler struct {
	animalUsecase usecase.AnimalUsecaseItf
}

func RegisterAnimalHandler(animalUsecase usecase.AnimalUsecaseItf, router fiber.Router) {
	animalHandler := animalHandler{animalUsecase}

	router = router.Group("/animals")
	router.Post("/_whatIs", middleware.BearerAuth, animalHandler.whatIs)
}

func (h *animalHandler) whatIs(c *fiber.Ctx) error {
	var raw model.PredictAnimalRequest
	if err := c.BodyParser(&raw); err != nil {
		return err
	}
	picture, err := c.FormFile("picture")
	if err != nil {
		return err
	}
	raw.Picture = picture
	resp, err := h.animalUsecase.PredictAnimal(c.Context(), &raw)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
