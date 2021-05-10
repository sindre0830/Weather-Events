package dict

// URL stores valid URL
var URL = "localhost"
var WEATHER_PATH = URL + "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = URL + "/weather-rest/v1/weather/compare/"
var WEATHEREVENT_PATH = URL + "/weather-rest/v1/event/date/"
var EVENT_PATH = URL + "/weather-rest/v1/weather/event/"
var DIAG_PATH = URL + "/weather-rest/v1/weather/diag/"
var HOLIDAY_PATH = URL + "/weather-rest/v1/notification/weather/holiday"
var WEATHERHOOK_PATH = URL + "/weather-rest/v1/notification/weather"

// Secret used for hashing
var Secret = []byte{43, 123, 65, 232, 4, 42, 35, 234, 21, 122, 214}

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
