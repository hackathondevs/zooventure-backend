package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/delivery/middleware"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type userHandler struct {
	usecase   usecase.UserUsecaseItf
	validator *validator.Validate
}

func RegisterUserHandler(
	usecase usecase.UserUsecaseItf,
	validator *validator.Validate,
	router fiber.Router,
) {
	userHandler := userHandler{usecase, validator}
	router = router.Group("/users")
	router.Get("", middleware.BearerAuth("false"), userHandler.profile)
	router.Put("", middleware.BearerAuth("false"), userHandler.updateProfile)
	router.Patch("", middleware.BearerAuth("false"), userHandler.resetPassword)
	router.Post("/_uploadProfilePicture", middleware.BearerAuth("false"), userHandler.updateProfilePicture)
	router.Delete("/_deleteProfilePicture", middleware.BearerAuth("false"), userHandler.deleteProfilePicture)
	router.Post("_exchange", middleware.BearerAuth("false"), userHandler.exchange)
	router.Get("_exchanges", middleware.BearerAuth("false"), userHandler.getExchanges)
}

func (h *userHandler) profile(c *fiber.Ctx) error {
	user, err := h.usecase.GetUserProfile(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (h *userHandler) updateProfile(c *fiber.Ctx) error {
	var profile model.UserCleanResource
	if err := c.BodyParser(&profile); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}
	if err := h.validator.Struct(&profile); err != nil {
		return err
	}
	if err := h.usecase.UpdateUserProfile(c.Context(), &profile); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *userHandler) resetPassword(c *fiber.Ctx) error {
	var attempt model.ResetPasswordRequest
	if err := c.BodyParser(&attempt); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}
	if err := h.validator.Struct(&attempt); err != nil {
		return err
	}
	if err := h.usecase.ResetPassword(c.Context(), attempt.CurrentPassword, attempt.Password); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *userHandler) updateProfilePicture(c *fiber.Ctx) error {
	picture, err := c.FormFile("profilePicture")
	if err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}
	if err := h.usecase.ChangePicture(c.Context(), picture); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *userHandler) deleteProfilePicture(c *fiber.Ctx) error {
	if err := h.usecase.DeletePicture(c.Context()); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *userHandler) exchange(c *fiber.Ctx) error {
	var exchangeReq model.ExchangeRequest
	if err := c.BodyParser(&exchangeReq); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}
	if err := h.validator.Struct(&exchangeReq); err != nil {
		return err
	}
	if err := h.usecase.Exchange(c.Context(), exchangeReq); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Exchange successful"})
}

func (h *userHandler) getExchanges(c *fiber.Ctx) error {
	exchanges, err := h.usecase.GetExchanges(c.Context())
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(exchanges)
}
