package usecase

import "errors"

var (
	ErrUserNotActive        = errors.New("User is not active")
	ErrUserNotExist         = errors.New("User is not exist")
	ErrWrongPassword        = errors.New("Wrong user password")
	ErrVerificationNotExist = errors.New("Verification not exist")
	ErrEmailExist           = errors.New("Email already exist")
	ErrNameExist            = errors.New("Name already exist")
	ErrIDNotNumeric         = errors.New("ID is not numeric")
	ErrInsufficientBalance  = errors.New("Insufficient balance")
	ErrMerchantNotExist     = errors.New("Merchant not exist")
	ErrCampaignNotExist     = errors.New("Campaign not exist")
)
