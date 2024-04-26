package usecase

import "errors"

var (
	ErrUserNotActive        = errors.New("El usuario no está activo/activa")
	ErrUserNotExist         = errors.New("El/La usuario no existe")
	ErrWrongPassword        = errors.New("La contraseña es incorrecta")
	ErrVerificationNotExist = errors.New("el intento de verificación no existe")
	ErrEmailExist           = errors.New("Ya existe el correo electrónico")
	ErrNameExist            = errors.New("Nombre ya existe")
	ErrIDNotNumeric         = errors.New("El valor de identificación no es válido")
)
