package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "harus ada":
		return ""
	case "ascii":
		return "harus berupa ascii"
	case "alphanum":
		return "harus berupa alfabet dan angka"
	case "alphanumunicode":
		return "harus berupa unicode"
	case "email":
		return "harus berupa email yang valid"
	case "min":
		return fmt.Sprintf("minimal karakter adalah %s", err.Param())
	case "eqfield":
		return fmt.Sprintf("hharus sama dengan %s", err.Param())
	case "url":
		return "harus berupa URL"
	default:
		return ""
	}
}
