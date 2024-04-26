package rest

import "errors"

var (
	ErrMissingID        = errors.New("falta identificación de verificación")
	ErrMissingToken     = errors.New("token de verificación faltante")
	ErrRequestMalformed = errors.New("solicitud mal formada")
)
