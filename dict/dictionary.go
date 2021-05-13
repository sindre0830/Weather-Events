package dict

// Secret used for hashing.
var Secret = []byte{43, 123, 65, 232, 4, 42, 35, 234, 21, 122, 214}

/* program paths */
var WEATHER_PATH = MAIN_URL + "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = MAIN_URL + "/weather-rest/v1/weather/compare/"
var WEATHEREVENT_PATH = MAIN_URL + "/weather-rest/v1/event/date/"
var EVENT_PATH = MAIN_URL + "/weather-rest/v1/weather/event/"
var DIAG_PATH = MAIN_URL + "/weather-rest/v1/weather/diag/"
var WEATHERHOOK_PATH = MAIN_URL + "/weather-rest/v1/notification/weather"

/* REST services */
var MAIN_URL = "localhost"
var YR_URL = "https://api.met.no/weatherapi/locationforecast/2.0/complete"

// GetYrURL generates yr REST service URL according to parameters.
func GetYrURL(latitude string, longitude string) string {
	return YR_URL + "?lat=" + latitude + "&lon=" + longitude
}

// GetWeatherURL generates Weather URL according to parameters.
func GetWeatherURL(location string, date string) string {
	if date != "" {
		return WEATHER_PATH + location + "?date=" + date
	}
	return WEATHER_PATH + location
}

// GetWeatherCompareURL generates WeatherCompare URL according to parameters.
func GetWeatherCompareURL(location string, date string) string {
	if date != "" {
		return WEATHERCOMPARE_PATH + location + "?date=" + date
	}
	return WEATHERCOMPARE_PATH + location
}
