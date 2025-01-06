package libs

import (
	"github.com/codingsince1985/geo-golang/openstreetmap"
	"log"
)

func FindCity(lat, lng float64) (string, error) {
	geo := openstreetmap.Geocoder()
	res, err := geo.ReverseGeocode(lat, lng)
	if err != nil {
		log.Printf("findCity: Cannot get City because: %v", err)
		return "", err
	}
	return res.City, nil
}
