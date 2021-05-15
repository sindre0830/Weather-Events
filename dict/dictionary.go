package dict

import "sync"

// We use mutex locks in callHook to ensure we get no concurrency issues WRT map writes
// Since we're loading every hook in when the program starts running, any 2 or more hooks with the same timeout
// would otherwise be at risk of panicking when running at the same time.
var MutexState = &sync.Mutex{}

// Secret used for hashing.
var Secret = []byte{43, 123, 65, 232, 4, 42, 35, 234, 21, 122, 214}

/* program paths */
var WEATHERDETAILS_PATH = "/weather-rest/v1/weather/location/"
var WEATHERCOMPARE_PATH = "/weather-rest/v1/weather/compare/"
var DIAG_PATH = "/weather-rest/v1/weather/diag/"
var WEATHER_HOOK_PATH = "/weather-rest/v1/notification/weather/"
var WEATHEREVENT_HOOK_PATH = "/weather-rest/v1/notification/event/"

/* REST services */
var MAIN_URL = "http://10.212.142.102"
var YR_URL = "https://api.met.no/weatherapi/locationforecast/2.0/complete"
var TICKETMASTER_URL = "https://app.ticketmaster.com/discovery/v2/events/"
var LOCATIONIQ_URL = "https://us1.locationiq.com/v1/search.php"
var PUBLICHOLIDAYS_URL = "https://date.nager.at/api/v2/PublicHolidays/"

/* API keys */
var TICKETMASTER_PK = ".json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot"
var LOCATIONIQ_PK = "?key=pk.d8a67c78822d16869c7a3e8f6d7617af"

/* collection names */
var COUNTRY_COLLECTION = "country-information"
var EVENT_COLLECTION = "event-information"
var LOCATION_COLLECTION = "location-information"
var HOLIDAYS_COLLECTION = "holidays-information"
var WEATHERDATA_COLLECTION = "weather-data"
var WEATHEREVENT_COLLECTION = "weather-event-hooks"
var WEATHER_COLLECTION = "weather-hooks"

func GetPublicHolidaysURL(year string, country string) string {
	return PUBLICHOLIDAYS_URL + year + "/" + country
}

// GetYrURL generates locationiq REST service URL according to parameters.
func GetLocationiqURL(location string) string {
	return LOCATIONIQ_URL + LOCATIONIQ_PK + "&q=" + location + "&format=json"
}

// GetYrURL generates ticketmaster REST service URL according to parameters.
func GetTicketmasterURL(event string) string {
	return TICKETMASTER_URL + event + TICKETMASTER_PK
}

// GetYrURL generates yr REST service URL according to parameters.
func GetYrURL(latitude string, longitude string) string {
	return YR_URL + "?lat=" + latitude + "&lon=" + longitude
}

// GetWeatherDetailsURL generates WeatherDetails URL according to parameters.
func GetWeatherDetailsURL(location string, date string) string {
	if date != "" {
		return MAIN_URL + WEATHERDETAILS_PATH + location + "?date=" + date
	}
	return WEATHERDETAILS_PATH + location
}

// GetWeatherCompareURL generates WeatherCompare URL according to parameters.
func GetWeatherCompareURL(location string, date string) string {
	if date != "" {
		return MAIN_URL + WEATHERCOMPARE_PATH + location + "?date=" + date
	}
	return WEATHERCOMPARE_PATH + location
}
