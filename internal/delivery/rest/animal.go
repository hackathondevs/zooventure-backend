package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
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
	router.Get("", middleware.BearerAuth("false"), animalHandler.fecthAll)
	router.Post("/trivia", middleware.BearerAuth("false"), animalHandler.trivia)
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

func (h *animalHandler) fecthAll(c *fiber.Ctx) error {
	animals, err := h.animalUsecase.FetchAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(animals)
}

func (h *animalHandler) trivia(c *fiber.Ctx) error {
	var animal struct {
		Name string `json:"name" validate:"required"`
	}
	if err := c.BodyParser(&animal); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}
	trivia, err := h.animalUsecase.GetTrivia(c.Context(), animal.Name)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(trivia)
}
