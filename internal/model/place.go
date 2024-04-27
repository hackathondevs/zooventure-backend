package model

type Place struct {
	Name  string `db:"Name"`
	Type  string `db:"Type"`
	Coord Coord  `db:"Coordinate"`
}

type PlcaeResource struct {
	Name string  `json:"name"`
	Type string  `db:"Type"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (e *Place) Clean() PlcaeResource {
	return PlcaeResource{
		Name: e.Name,
		Type: e.Type,
		Lat:  e.Coord.Lat,
		Long: e.Coord.Long,
	}
}
