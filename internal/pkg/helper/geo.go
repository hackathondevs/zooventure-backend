package helper

import (
	"math"

	"github.com/mirzahilmi/hackathon/internal/model"
)

const earthRadius = 6371

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func Haversine(x, y model.Coord) float64 {
	aLatRad := degreesToRadians(x.Lat)
	aLonRad := degreesToRadians(x.Long)
	bLatRad := degreesToRadians(y.Lat)
	bLonRad := degreesToRadians(y.Long)
	deltaLat := bLatRad - aLatRad
	deltaLon := bLonRad - aLonRad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(aLatRad)*math.Cos(bLatRad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	return distance
}
