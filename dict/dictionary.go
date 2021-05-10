package dict

// URL stores valid URL
var URL = "localhost"
var WEATHER_PATH = URL + "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = URL + "/weather-rest/v1/weather/compare/"
var EVENT_PATH = URL + "/weather-rest/v1/weather/event/"
var DIAG_PATH = URL + "/weather-rest/v1/weather/diag/"

// GetWeatherURL generates Weather URL according to parameters
func GetWeatherURL(location string, date string) string {
	if date != "" {
		return WEATHER_PATH + location + "?date=" + date
	}
	return WEATHER_PATH + location
}

// GetWeatherCompareURL generates WeatherCompare URL according to parameters
func GetWeatherCompareURL(location string, date string) string {
	if date != "" {
		return WEATHERCOMPARE_PATH + location + "?date=" + date
	}
	return WEATHERCOMPARE_PATH + location
}
