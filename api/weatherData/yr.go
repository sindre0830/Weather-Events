package weatherData

type Yr struct {
	Properties struct {
		Meta       interface{}   `json:"meta"`
		Timeseries []interface{} `json:"timeseries"`
	} `json:"properties"`
}
