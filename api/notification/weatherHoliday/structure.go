package weatherHoliday

// WeatherHolidayInput structure, stores information from the user about the webhook
type WeatherHolidayInput struct {
	Holiday   string `json:"holiday"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}

// WeatherHoliday structure, stores information about the webhook added to the database
type WeatherHoliday struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Holiday   string `json:"holiday"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}
