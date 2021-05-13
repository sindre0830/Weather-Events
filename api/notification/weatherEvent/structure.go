package weatherEvent

type WeatherEvent struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}

type WeatherEventDefault struct {
	Date      string `json:"date"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}

type WeatherEventHoliday struct {
	Holiday   string `json:"holiday"`
	Location  string `json:"location"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}

type WeatherEventTicketmaster struct {
	Ticket    string `json:"ticket"`
	URL       string `json:"url"`
	Frequency string `json:"frequency"`
	Timeout   int64  `json:"timeout"`
}
