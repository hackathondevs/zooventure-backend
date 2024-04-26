package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/model"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

type authHandler struct {
	usecase   usecase.AuthUsecaseItf
	validator *validator.Validate
}

func RegisterAuthHandler(
	usecase usecase.AuthUsecaseItf,
	validator *validator.Validate,
	router fiber.Router,
) {
	authHandler := authHandler{usecase, validator}
	router = router.Group("/users")
	router.Post("", authHandler.signUp)
	router.Post("/_login", authHandler.login)
	router.Get("/_verify", authHandler.verify)
}

func (h *authHandler) signUp(c *fiber.Ctx) error {
	var payloadUser model.UserSignupRequest
	if err := c.BodyParser(&payloadUser); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}

	if err := h.validator.Struct(&payloadUser); err != nil {
		return err
	}

	if err := h.usecase.RegisterUser(c.Context(), payloadUser); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *authHandler) verify(c *fiber.Ctx) error {
	verifReq := model.UserVerifRequest{ID: c.Query("s"), Token: c.Query("t")}
	if err := h.validator.Struct(&verifReq); err != nil {
		return err
	}

	if err := h.usecase.VerifyUser(c.Context(), verifReq); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *authHandler) login(c *fiber.Ctx) error {
	var attempt model.UserLoginRequest
	if err := c.BodyParser(&attempt); err != nil {
		return response.NewError(fiber.StatusBadRequest, ErrRequestMalformed)
	}

	if err := h.validator.Struct(&attempt); err != nil {
		return err
	}

	signedToken, err := h.usecase.LogUserIn(c.Context(), attempt)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": signedToken})
}
