package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
)

func RegisterAnimalHandler(router fiber.Router) {
	router = router.Group("/animals")
	router.Post("/_whatIs", whatIs)
}

func whatIs(c *fiber.Ctx) error {
	_, err := c.FormFile("picture")
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}

	// Call usecase to fetch picture context (prompt)

	/*
		{
			"name": "Harimau",
			"latin": "Harimau",
			"habitat": "Sidoarjo, Indonesia",
			"diets": "Karnivora",
			"lifespan": "35 Tahun",
			"funfact": "Berbeda dengan singa, harimau dewasa adalah satwa soliter yang menandai wilayahnya dengan urin dan cakaran di batang pohon. (max 20 kata)"
		}
	*/

	return c.Status(fiber.StatusOK).JSON(model.Animal{})
}
