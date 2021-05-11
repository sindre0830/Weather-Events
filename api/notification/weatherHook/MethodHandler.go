package weatherHook

import (
	"main/debug"
	"net/http"
)

// MethodHandler - handles the method of a http request.
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	var weatherHook WeatherHook
	switch r.Method {
	case http.MethodPost:
		weatherHook.HandlerPost(w, r)
	case http.MethodGet:
		weatherHook.HandlerGet(w, r)
	case http.MethodDelete:
		weatherHook.HandlerDelete(w, r)
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
