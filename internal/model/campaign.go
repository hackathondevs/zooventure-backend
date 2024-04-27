package model

import "mime/multipart"

type Campaign struct {
	ID          int64  `db:"ID" json:"id,omitempty"`
	Picture     string `db:"Picture" json:"picture"`
	Title       string `db:"Title" json:"title"`
	Description string `db:"Description" json:"description,omitempty"`
	Reward      int    `db:"Reward" json:"reward"`
	Submitted   bool   `db:"Submitted" json:"submitted"`
	Submission  string `db:"Submission" json:"submission,omitempty"`
}

type CampaignRequest struct {
	Picture     *multipart.FileHeader `form:"picture" validate:"required"`
	Title       string                `form:"title" validate:"required"`
	Description string                `form:"description" validate:"required"`
	Reward      int                   `form:"reward" validate:"required"`
}

type CampaignSubmissionRequest struct {
	Submission string `json:"submission" validate:"required"`
}

type CampaignSubmission struct {
	UserID     int64  `db:"UserID" json:"userId"`
	CampaignID int64  `db:"CampaignID" json:"campaignId"`
	Submission string `db:"Submission" json:"submission"`
}
