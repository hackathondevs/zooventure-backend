package rest

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type campaignHandler struct {
	usecase usecase.CampaignUsecaseItf
}

func RegisterCampaignHandler(
	usecase usecase.CampaignUsecaseItf,
	router fiber.Router,
) {
	campaignHandler := campaignHandler{usecase}
	router = router.Group("/campaigns")
	router.Get("", campaignHandler.FetchAll)
	router.Get("/:id<int>", campaignHandler.FetchSingle)
}

func (h *campaignHandler) FetchAll(c *fiber.Ctx) error {
	campaigns, err := h.usecase.FetchAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(campaigns)
}

func (h *campaignHandler) FetchSingle(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	campaign, err := h.usecase.GetWithSubmission(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(campaign)
}
