package dict

// URL stores valid URL
var URL = "localhost"
var WEATHER_PATH = URL + "/weather-rest/v1/weather/location/"

// Generate Weather URL according to parameters
func getWeatherURL(location string, date string) string {
	if date != "" {
		return WEATHER_PATH + location + "?date=" + date
	}
	return WEATHER_PATH + location
}
