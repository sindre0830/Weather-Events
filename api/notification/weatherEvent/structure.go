package weatherEvent

type WeatherEventInput struct {
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}

type WeatherEvent struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}
