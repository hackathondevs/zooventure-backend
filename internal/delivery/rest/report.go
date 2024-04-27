package rest

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecaseItf
	validator     *validator.Validate
}

func RegisterReportHandler(
	usecase usecase.ReportUsecaseItf,
	validator *validator.Validate,
	router fiber.Router,
) {
	reportHandler := ReportHandler{usecase, validator}
	router = router.Group("/reports")
	router.Post("", middleware.BearerAuth("false"), reportHandler.CreateReport)
	router.Get("", middleware.BearerAuth("true"), reportHandler.GetReports)
	router.Patch("/:id<int>", middleware.BearerAuth("true"), reportHandler.UpdateReport)
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error {
	pict, err := c.FormFile("picture")
	if err != nil {
		return err
	}

	req := model.ReportRequest{
		Picture:     pict,
		Description: c.FormValue("description"),
		Location:    c.FormValue("location"),
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}
	if err := h.reportUsecase.CreateReport(c.Context(), req); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Report created"})
}

func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	reports, err := h.reportUsecase.GetReports(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(reports)
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return err
	}

	value := c.Query("action")
	if err := h.reportUsecase.UpdateReport(c.Context(), id, value); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Report updated"})
}
