package model

type Campaign struct {
	ID          int64  `db:"ID" json:"id,omitempty"`
	Picture     string `db:"Picture" json:"picture"`
	Title       string `db:"Title" json:"title"`
	Description string `db:"Description" json:"description,omitempty"`
	Reward      int    `db:"Reward" json:"reward"`
	Submitted   bool   `db:"Submitted" json:"submitted"`
	Submission  string `db:"Submission" json:"submission,omitempty"`
}
