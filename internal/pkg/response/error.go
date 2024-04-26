package response

import (
	"errors"

	"github.com/gofiber/fiber/v2/utils"
)

type Error struct {
	Code int
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func NewError(code int, err error) error {
	return &Error{code, err}
}

func NewHTTPError(code int) error {
	return &Error{code, errors.New(utils.StatusMessage(code))}
}
