package model

import "mime/multipart"

type ReportResource struct {
	ID          int64  `json:"id" db:"ID"`
	Picture     string `json:"picture" db:"Picture"`
	Description string `json:"description" db:"Description"`
	Location    string `json:"location" db:"Location"`
	CreatedAt   string `json:"createdAt" db:"CreatedAt"`
	Action      string `json:"action" db:"Action"`
}

type ReportRequest struct {
	Picture     *multipart.FileHeader `form:"picture" validate:"required"`
	Description string                `form:"description" validate:"required"`
	Location    string                `form:"location" validate:"required"`
}
