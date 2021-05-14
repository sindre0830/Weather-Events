package weather

import (
	"main/debug"
	"net/http"
)

// MethodHandler handles the method of a http request.
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	var weather Weather
	switch r.Method {
		case http.MethodGet:
			weather.get(w, r)
		case http.MethodPost:
			weather.post(w, r)
		case http.MethodDelete:
			weather.delete(w, r)
		default:
			debug.ErrorMessage.Update(
				http.StatusMethodNotAllowed,
				"weather.MethodHandler() -> Validating method",
				"Method validation: wrong method",
				"Method not implemented.",
			)
			debug.ErrorMessage.Print(w)
	}
}
