package model

type Trivia struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	OptA     string `json:"optA"`
	OptB     string `json:"optB"`
	OptC     string `json:"optC"`
}
