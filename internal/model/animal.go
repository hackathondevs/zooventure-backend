package model

import "mime/multipart"

type Animal struct {
	Name            string   `json:"name"`
	Latin           string   `json:"latin"`
	CountryOfOrigin string   `json:"countryOfOrigin"`
	Characteristics []string `json:"characteristics"`
	Category        string   `json:"category"`
	Lifespan        string   `json:"lifespan"`
	Funfact         string   `json:"funfact"`
	GotBonus        bool     `json:"gotBonus"`
}

type PredictAnimalRequest struct {
	Picture *multipart.FileHeader `form:"-"`
	Lat     float64               `form:"lat" validate:"required"`
	Long    float64               `form:"long" validate:"required"`
}
