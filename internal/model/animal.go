package model

import "mime/multipart"

type Animal struct {
	ID                  int64   `db:"ID"`
	Picture             string  `db:"Picture"`
	Name                string  `db:"Name"`
	Latin               string  `db:"Latin"`
	Origin              string  `db:"Origin"`
	Characteristic      string  `db:"Characteristic"`
	Diet                string  `db:"Diet"`
	Lifespan            string  `db:"Lifespan"`
	EnclosureCoordinate Coord   `db:"EnclosureCoordinate"`
	Distance            float64 `db:"Distance"`
}

type AnimalResource struct {
	ID             int64   `json:"id"`
	Picture        string  `json:"picture"`
	Name           string  `json:"name"`
	Latin          string  `json:"latin"`
	Origin         string  `json:"origin"`
	Characteristic string  `json:"characteristic"`
	Diet           string  `json:"diet"`
	Lifespan       string  `json:"lifespan"`
	Lat            float64 `json:"lat"`
	Long           float64 `json:"long"`
}

type PredictAnimalReq struct {
	Picture *multipart.FileHeader `form:"picture"`
	Lat     float64               `form:"lat"`
	Long    float64               `form:"long"`
}

func (a *Animal) Resource() AnimalResource {
	return AnimalResource{
		ID:             a.ID,
		Picture:        a.Picture,
		Name:           a.Name,
		Latin:          a.Latin,
		Origin:         a.Origin,
		Characteristic: a.Characteristic,
		Diet:           a.Diet,
		Lifespan:       a.Lifespan,
		Lat:            a.EnclosureCoordinate.Lat,
		Long:           a.EnclosureCoordinate.Long,
	}
}
