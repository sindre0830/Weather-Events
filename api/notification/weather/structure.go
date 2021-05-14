package weather

// Our main struct with ID
type Weather struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	URL      string `json:"url"`
	Timeout  int64  `json:"timeout"`
}

// Helper struct for reading from user
type WeatherInput struct {
	Location string `json:"location"`
	URL      string `json:"url"`
	Timeout  int64  `json:"timeout"`
}
