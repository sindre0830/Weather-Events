package geoCoords

/**
*	LocationCoords
*	Holds our latitude and longitude data for one location
**/
type LocationCoords struct {
	Address    string  `json:"address"`
	Importance float64 `json:"importance"`
	Latitude   float64 `json:"lat"`
	Longitude  float64 `json:"lon"`
}
