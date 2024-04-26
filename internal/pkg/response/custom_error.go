package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type CustomError struct {
	Code   int
	Errors fiber.Map
}

func (e *CustomError) Error() string {
	return utils.StatusMessage(e.Code)
}

func NewCustomError(code int, errs fiber.Map) error {
	return &CustomError{code, errs}
}
