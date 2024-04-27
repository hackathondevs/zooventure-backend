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
		return "Harus berupa ascii"
	case "alphanum":
		return "Harus berupa alfabet dan angka"
	case "alphanumunicode":
		return "Harus berupa unicode"
	case "email":
		return "Harus berupa email yang valid"
	case "min":
		return fmt.Sprintf("Minimal memiliki panjang %s karakter", err.Param())
	case "eqfield":
		return fmt.Sprintf("Harus sama dengan %s", err.Param())
	case "url":
		return "Harus berupa URL"
	default:
		return ""
	}
}
