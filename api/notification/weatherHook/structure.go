package weatherHook

type WeatherHook struct {
	ID       string `json:"id"`
	Location string `json:"location"`
	Timeout  int64  `json:"timeout"`
	URL      string `json:"url"`
}

type WeatherHookInput struct {
	Location string `json:"location"`
	Timeout  int64  `json:"timeout"`
	URL      string `json:"url"`
}
