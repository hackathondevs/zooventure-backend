package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type animalHandler struct {
	animalUsecase usecase.AnimalUsecaseItf
}

func RegisterAnimalHandler(animalUsecase usecase.AnimalUsecaseItf, router fiber.Router) {
	animalHandler := animalHandler{animalUsecase}

	router = router.Group("/animals")
	router.Post("/_whatIs", animalHandler.whatIs)
}

func (h *animalHandler) whatIs(c *fiber.Ctx) error {
	pict, err := c.FormFile("picture")
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}

	resp := h.animalUsecase.PredictAnimal(c.Context(), pict)

	return c.Status(fiber.StatusOK).JSON(resp)
}
