package helper

import (
	"github.com/go-playground/validator/v10"
)

func ValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return ""
	case "ascii":
		return ""
	case "alphanum":
		return ""
	case "alphanumunicode":
		return ""
	case "e164":
		return ""
	case "email":
		return ""
	case "min":
		return ""
	case "eqfield":
		return ""
	case "url":
		return ""
	default:
		return ""
	}
}
