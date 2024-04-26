package model

type Animal struct {
	Name     string `json:"name"`
	Latin    string `json:"latin"`
	Habitat  string `json:"habitat"`
	Diets    string `json:"diets"`
	Lifespan string `json:"lifespan"`
	Funfact  string `json:"funfact"`
}
