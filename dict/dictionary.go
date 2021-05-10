package dict

// URL stores valid URL
var URL = "localhost"
var WEATHER_PATH = URL + "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = URL + "/weather-rest/v1/weather/compare/"
var WEATHEREVENT_PATH = URL + "/weather-rest/v1/event/date/"
var EVENT_PATH = URL + "/weather-rest/v1/weather/event/"
var HOLIDAY_PATH = URL + "/weather-rest/v1/notification/weather/holiday/"

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
