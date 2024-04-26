package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mirzahilmi/hackathon/internal/pkg/pasetok"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/mirzahilmi/hackathon/internal/usecase"
)

func BearerAuth(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return response.NewHTTPError(fiber.StatusUnauthorized)
	}
	bearer := strings.SplitN(header, " ", 2)
	if bearer[0] != "Bearer" || len(bearer) != 2 {
		return response.NewHTTPError(fiber.StatusBadRequest)
	}
	token, err := pasetok.Decode(bearer[1])
	if err != nil {
		return response.NewHTTPError(fiber.StatusUnauthorized)
	}
	sub, err := token.GetSubject()
	if err != nil {
		return err
	}
	name, err := token.GetString("name")
	if err != nil {
		return err
	}
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return err
	}
	c.Locals(usecase.ClientID, id)
	c.Locals(usecase.ClientName, name)
	return c.Next()
}
