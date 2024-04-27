package rest

import "errors"

var (
	ErrMissingID        = errors.New("kode identitas tidak dapat ditemukan")
	ErrMissingToken     = errors.New("token tidak dapat ditemukan")
	ErrRequestMalformed = errors.New("Bentuk request tidak dapat diproses")
)
