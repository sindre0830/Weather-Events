package weatherHook

type WeatherHook struct {
	Location string `json:"location"`
	Timeout  int64  `json:"timeout"`
}
