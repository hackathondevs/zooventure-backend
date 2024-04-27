package rest

import (
	"mime/multipart"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type campaignHandler struct {
	usecase   usecase.CampaignUsecaseItf
	validator *validator.Validate
}

func RegisterCampaignHandler(
	usecase usecase.CampaignUsecaseItf,
	router fiber.Router,
	validator *validator.Validate,
) {
	campaignHandler := campaignHandler{usecase, validator}
	router = router.Group("/campaigns")
	router.Get("", middleware.BearerAuth("false"), campaignHandler.FetchAll)
	router.Get("/:id<int>", middleware.BearerAuth("false"), campaignHandler.FetchSingle)
	router.Post("/_submission/:id<int>", middleware.BearerAuth("false"), campaignHandler.SubmitSubmission)
	router.Post("", middleware.BearerAuth("true"), campaignHandler.Create)
	router.Patch("/:id<int>", middleware.BearerAuth("true"), campaignHandler.Update)
	router.Delete("/:id<int>", middleware.BearerAuth("true"), campaignHandler.Delete)
	router.Get("/_admin", middleware.BearerAuth("true"), campaignHandler.GetAllForAdmin)
	router.Get("/_admin/:id<int>", middleware.BearerAuth("true"), campaignHandler.GetByIDForAdmin)
	router.Get("/_submissions", middleware.BearerAuth("true"), campaignHandler.GetAllCampaignSubmission)
	router.Patch("/_submissions/:id<int>/:userID<int>", middleware.BearerAuth("true"), campaignHandler.UpdateStatusCampaignSubmission)
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

func (h *campaignHandler) SubmitSubmission(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	var req model.CampaignSubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	err = h.usecase.SubmitSubmission(c.Context(), id, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Submission uploaded successfully"})
}

func (h *campaignHandler) Create(c *fiber.Ctx) error {
	pict, err := c.FormFile("picture")
	if err != nil {
		return err
	}

	reward, err := strconv.ParseUint(c.FormValue("reward"), 10, 64)
	if err != nil {
		return err
	}

	req := model.CampaignRequest{
		Picture:     pict,
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Reward:      int(reward),
	}

	if err := h.validator.Struct(&req); err != nil {
		return err
	}

	err = h.usecase.Create(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Campaign created successfully"})
}

func (h *campaignHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	var pict *multipart.FileHeader
	pict, err = c.FormFile("picture")
	if err != nil {
		if err.Error() == "there is no uploaded file associated with the given key" {
			pict = nil
		} else {
			return err
		}
	}

	var reward uint64
	if c.FormValue("reward") == "" {
		reward = 0
	} else {
		reward, err = strconv.ParseUint(c.FormValue("reward"), 10, 64)
		if err != nil {
			return err
		}
	}

	req := model.CampaignRequest{
		Picture:     pict,
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Reward:      int(reward),
	}

	err = h.usecase.Update(c.Context(), req, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Campaign updated successfully"})

}

func (h *campaignHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	err = h.usecase.Delete(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Campaign deleted successfully"})
}

func (h *campaignHandler) GetAllForAdmin(c *fiber.Ctx) error {
	campaigns, err := h.usecase.GetAll(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(campaigns)
}

func (h *campaignHandler) GetByIDForAdmin(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	campaign, err := h.usecase.GetByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(campaign)
}

func (h *campaignHandler) GetAllCampaignSubmission(c *fiber.Ctx) error {
	campaigns, err := h.usecase.GetAllCampaignSubmission(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(campaigns)
}

func (h *campaignHandler) UpdateStatusCampaignSubmission(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	userID, err := strconv.ParseInt(c.Params("userID"), 10, 64)
	if err != nil {
		return err
	}

	value := c.Query("status")
	err = h.usecase.UpdateStatusCampaignSubmission(c.Context(), id, userID, value)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Submission updated"})
}
