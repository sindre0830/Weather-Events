package weatherDetails

import (
	"main/debug"
	"net/http"
)

// MethodHandler handles the method of a http request.
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case http.MethodGet:
			var weatherDetails WeatherDetails
			weatherDetails.Handler(w, r)
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
