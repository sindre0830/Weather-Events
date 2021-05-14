package dict

import "sync"

// We use mutex locks in callHook to ensure we get no concurrency issues WRT map writes
// Since we're loading every hook in when the program starts running, any 2 or more hooks with the same timeout
// would otherwise be at risk of panicking when running at the same time.
var MutexState = &sync.Mutex{}

// Secret used for hashing.
var Secret = []byte{43, 123, 65, 232, 4, 42, 35, 234, 21, 122, 214}

/* program paths */
var WEATHERDETAILS_PATH = MAIN_URL + "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = MAIN_URL + "/weather-rest/v1/weather/compare/"
var DIAG_PATH = MAIN_URL + "/weather-rest/v1/weather/diag/"
var WEATHER_HOOK_PATH = MAIN_URL + "/weather-rest/v1/notification/weather/"
var WEATHEREVENT_HOOK_PATH = MAIN_URL + "/weather-rest/v1/notification/event/"

/* REST services */
var MAIN_URL = "localhost"
var YR_URL = "https://api.met.no/weatherapi/locationforecast/2.0/complete"

// GetYrURL generates yr REST service URL according to parameters.
func GetYrURL(latitude string, longitude string) string {
	return YR_URL + "?lat=" + latitude + "&lon=" + longitude
}

// GetWeatherDetailsURL generates WeatherDetails URL according to parameters.
func GetWeatherDetailsURL(location string, date string) string {
	if date != "" {
		return WEATHERDETAILS_PATH + location + "?date=" + date
	}
	return WEATHERDETAILS_PATH + location
}

// GetWeatherCompareURL generates WeatherCompare URL according to parameters.
func GetWeatherCompareURL(location string, date string) string {
	if date != "" {
		return WEATHERCOMPARE_PATH + location + "?date=" + date
	}
	return WEATHERCOMPARE_PATH + location
}
