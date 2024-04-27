package model

type Enclosure struct {
	Name  string `db:"Name"`
	Coord Coord  `db:"Coordinate"`
}

type EnclosureResource struct {
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (e *Enclosure) Clean() EnclosureResource {
	return EnclosureResource{
		Name: e.Name,
		Lat:  e.Coord.Lat,
		Long: e.Coord.Long,
	}
}
