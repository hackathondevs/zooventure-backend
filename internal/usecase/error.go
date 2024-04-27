package usecase

import "errors"

var (
	ErrUserNotActive        = errors.New("Pengguna masih belum terverifikasi")
	ErrUserNotExist         = errors.New("Pengguna tidak dapat ditemukan")
	ErrWrongPassword        = errors.New("Kata sandi salah")
	ErrVerificationNotExist = errors.New("Percobaan verifikasi tidak ditemukan")
	ErrEmailExist           = errors.New("Alamat email sudah digunakan")
	ErrIDNotNumeric         = errors.New("Nomor ID bukan berupa angka")
	ErrInsufficientBalance  = errors.New("Saldo tidak cukup")
	ErrMerchantNotExist     = errors.New("Pedagang tidak dapat ditemukan")
	ErrCampaignNotExist     = errors.New("Tidak dapat menemukan campaign")
	ErrAlreadyReported      = errors.New("Sudah melaporkan campaign ini")
)
