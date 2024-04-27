package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type animalHandler struct {
	animalUsecase usecase.AnimalUsecaseItf
	validator     *validator.Validate
}

func RegisterAnimalHandler(
	animalUsecase usecase.AnimalUsecaseItf,
	validator *validator.Validate,
	router fiber.Router,
) {
	animalHandler := animalHandler{animalUsecase, validator}

	router = router.Group("/animals")
	router.Post("/_whatIs", middleware.BearerAuth("false"), animalHandler.whatIs)
}

func (h *animalHandler) whatIs(c *fiber.Ctx) error {
	var raw model.PredictAnimalReq
	if err := c.BodyParser(&raw); err != nil {
		return err
	}
	if err := h.validator.Struct(&raw); err != nil {
		return err
	}
	picture, err := c.FormFile("picture")
	if err != nil {
		return err
	}
	raw.Picture = picture
	animal, err := h.animalUsecase.PredictAnimal(c.Context(), &raw)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(animal)
}
