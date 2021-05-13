package weatherEvent

import (
	"main/debug"
	"net/http"
)

// MethodHandler handles the method of a http request.
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	var weatherEvent WeatherEvent
	switch r.Method {
		case http.MethodGet:
			weatherEvent.get(w, r)
		case http.MethodPost:
			weatherEvent.post(w, r)
		case http.MethodDelete:
			weatherEvent.delete(w, r)
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
