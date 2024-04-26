package model

type UserVerificationResource struct {
	ID     int64  `db:"ID"`
	UserID int64  `db:"UserID"`
	Token  string `db:"Token"`
}

type UserVerifRequest struct {
	ID    string `json:"id" validate:"required,numeric"`
	Token string `json:"token" validate:"required,alphanum,len=32"`
}
