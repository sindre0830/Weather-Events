package diag

import (
	"encoding/json"
	"log"
	"main/debug"
	"math"
	"net/http"
	"time"
)

var StartTime time.Time
var HookAmount int

//MethodHandler -Class function will be called and handle all requests and fetches
func MethodHandler(w http.ResponseWriter, r *http.Request) {
	var diag DiagStatuses
	var err error

	diag.Restcountries, err = getStatusOf("https://restcountries.eu/rest/v2/all")
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Diag.Handler() -> Diag.getStatusOf() -> getting status of APIs",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	diag.TicketMaster, err = getStatusOf("https://app.ticketmaster.com/discovery/v2/events.json?apikey=ySyIqc6FFKgUIIgzKB5LAOQeGUeU1mot")
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Diag.Handler() -> Diag.getStatusOf() -> getting status of APIs",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	diag.LocationIq, err = getStatusOf("https://eu1.locationiq.com/v1/reverse.php?key=pk.d8a67c78822d16869c7a3e8f6d7617af&lat=32&lon=60&format=json")
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Diag.Handler() -> Diag.getStatusOf() -> getting status of APIs",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//Note that this one will almost always return 403, see their documentation here: https://api.met.no/doc/FAQ
	diag.Weatherapi, err = getStatusOf("https://api.met.no/weatherapi/locationforecast/2.0/complete?lat=30.0&lon=30.0")
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Diag.Handler() -> Diag.getStatusOf() -> getting status of APIs",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}
	//Note that this one will almost always return 405, no documentation was found on this, however the endpoints using this api works fine.
	diag.PublicHolidays, err = getStatusOf("https://date.nager.at/api/v2/publicholidays/2021/NO")
	if err != nil {
		debug.ErrorMessage.Update(
			http.StatusInternalServerError,
			"Diag.Handler() -> Diag.getStatusOf() -> getting status of APIs",
			err.Error(),
			"Unknown",
		)
		debug.ErrorMessage.Print(w)
		return
	}

	diag.RegisteredWebhooks = HookAmount
	diag.Version = "v1"
	diag.Uptime = int(math.Floor(time.Since(StartTime).Seconds()))

	//Formats the printouts
	w.Header().Set("Content-Type", "application/json")
	//Outputs results
	err = json.NewEncoder(w).Encode(diag)
	if err != nil {
		log.Println("ERROR encoding JSON", err) //If error, send log to error logger
	}
}

//getStatusOf returns the status code of a head request to the root path of a remote.
func getStatusOf(addr string) (int, error) {
	res, err := http.Head(addr)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	res.Body.Close()
	return res.StatusCode, nil
}
