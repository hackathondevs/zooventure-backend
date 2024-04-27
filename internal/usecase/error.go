package usecase

import "errors"

var (
	ErrUserNotActive        = errors.New("Pengguna belum aktif")
	ErrUserNotExist         = errors.New("Pengguna tidak ditemukan")
	ErrWrongPassword        = errors.New("Kata sandi salah")
	ErrVerificationNotExist = errors.New("Percobaan verifikasi tidak ditemukan")
	ErrEmailExist           = errors.New("Email telah digunakan")
	ErrNameExist            = errors.New("Nama telah digunakan")
	ErrIDNotNumeric         = errors.New("kode id bukan angka")
)
