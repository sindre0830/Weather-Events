package holidaysData

// Struct for information about one holiday, used when getting data from the API
type Holiday []struct {
	Date string `json:"date"`
	Name string `json:"name"`
}
