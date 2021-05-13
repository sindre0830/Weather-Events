package weatherHook

// Our main struct with ID
type WeatherHook struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Timeout  int64  `json:"timeout"`
	URL      string `json:"url"`
}

// Helper struct for reading from user
type WeatherHookInput struct {
	Location string `json:"location"`
	Timeout  int64  `json:"timeout"`
	URL      string `json:"url"`
}
