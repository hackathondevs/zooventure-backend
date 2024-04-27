package conf

import (
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	jsoniter "github.com/json-iterator/go"
	"github.com/mirzahilmi/hackathon/internal/pkg/response"
	"github.com/sirupsen/logrus"
)

var StartTime = time.Now()

func NewFiber(log *logrus.Logger) *fiber.App {
	fiber := fiber.New(fiber.Config{
		AppName:           os.Getenv("APP_NAME"),
		Prefork:           true,
		BodyLimit:         50 * 1024 * 1024,
		DisableKeepalive:  true,
		StrictRouting:     true,
		CaseSensitive:     true,
		EnablePrintRoutes: true,
		ErrorHandler:      newErrorHandler(log),
		JSONEncoder:       jsoniter.Marshal,
		JSONDecoder:       jsoniter.Unmarshal,
	})
	return fiber
}

func newErrorHandler(log *logrus.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		log.Error(err)
		var apiErr *response.Error
		if errors.As(err, &apiErr) {
			return ctx.Status(apiErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": apiErr.Error()},
			})
		}

		if validationErr, ok := err.(validator.ValidationErrors); ok {
			fieldErr := fiber.Map{}
			for _, e := range validationErr {
				fieldErr[e.Field()] = e.Error()
			}
			fieldErr["message"] = utils.StatusMessage(fiber.StatusUnprocessableEntity)
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": fieldErr,
			})
		}

		var apiCustomErr *response.CustomError
		if errors.As(err, &apiCustomErr) {
			return ctx.Status(apiCustomErr.Code).JSON(fiber.Map{
				"errors": apiCustomErr.Errors,
			})
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return ctx.Status(fiberErr.Code).JSON(fiber.Map{
				"errors": fiber.Map{"message": utils.StatusMessage(fiberErr.Code)},
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fiber.Map{"message": utils.StatusMessage(fiber.StatusInternalServerError)},
		})
	}
}
