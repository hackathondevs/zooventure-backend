package model

import (
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

type UserResource struct {
	ID             int64  `db:"ID"`
	Email          string `db:"Email"`
	Password       string `db:"Password"`
	ProfilePicture string `db:"ProfilePicture"`
	Name           string `db:"Name"`
	Active         bool   `db:"Active"`
	Admin          bool   `db:"Admin"`
	Balance        int    `db:"Balance"`
}

type UserCleanResource struct {
	Email          string `json:"email" validate:"-"`
	ProfilePicture string `json:"profilePicture" validate:"-"`
	Name           string `json:"name" validate:"required,alphanumunicode"`
	Balance        int    `json:"balance" validate:"-"`
}

type UserSignupRequest struct {
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
	Name                 string `json:"name" validate:"required,ascii"`
}

type UserLoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8"`
	RememberMe bool   `json:"rememberMe" validate:"boolean"`
}

type ResetPasswordRequest struct {
	CurrentPassword      string `json:"currentPassword" validate:"required,min=8"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"required,eqfield=Password"`
}

func (u UserResource) Clean(supabase *storage_go.Client) (UserCleanResource, error) {
	if u.ProfilePicture != "" {
		u.ProfilePicture = supabase.GetPublicUrl(os.Getenv("SUPABASE_BUCKET_ID"), u.ProfilePicture).SignedURL
	}
	return UserCleanResource{
		Email:          u.Email,
		ProfilePicture: u.ProfilePicture,
		Name:           u.Name,
		Balance:        u.Balance,
	}, nil
}

type ExchangeRequest struct {
	Amount int64  `json:"amount" validate:"required,min=1"`
	Code   string `json:"code" validate:"required"`
}

type ExchangeResource struct {
	ID           int64  `db:"ID"`
	UserID       int64  `db:"UserID"`
	MerchantID   int64  `db:"MerchantID"`
	MerchantName string `db:"MerchantName"`
	MerchantCode string `db:"MerchantCode"`
	Amount       int64  `db:"Amount"`
	Status       string `db:"Status"`
	Date         string `db:"Date"`
}

type ExchangeCleanResource struct {
	MerchantName string `json:"merchantName"`
	MerchantCode string `json:"merchantCode"`
	Amount       int64  `json:"amount"`
	Status       string `json:"status"`
	Date         string `json:"date"`
}

func (e ExchangeResource) Clean() ExchangeCleanResource {
	return ExchangeCleanResource{
		MerchantName: e.MerchantName,
		MerchantCode: e.MerchantCode,
		Amount:       e.Amount,
		Status:       e.Status,
		Date:         e.Date,
	}
}
