package repository

import "errors"

var (
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrNotTransaction = errors.New("not transaction")
)
