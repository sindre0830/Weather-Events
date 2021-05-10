package weatherHoliday

import (
	"main/debug"
	"net/http"
)

// MethodHandler handles the method of a http request.
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	var weatherHoliday WeatherHoliday

	switch r.Method {
	case http.MethodPost:
		weatherHoliday.POST(w, r)
	case http.MethodDelete:
		weatherHoliday.Delete(w, r)
	default:
		debug.ErrorMessage.Update(
			http.StatusMethodNotAllowed,
			"MethodHandler() -> Validating method",
			"method validation: wrong method",
			"Method not implemented.",
		)
		debug.ErrorMessage.Print(w)
	}
}
