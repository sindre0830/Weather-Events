package debug

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ErrorMessage is used for all error handling in the program.
var ErrorMessage Debug

// Debug structure stores information about errors.
//
// Functionality: Update, Print
type Debug struct {
	StatusCode 		 int    `json:"status_code"`
	Location   		 string `json:"location"`
	RawError   		 string `json:"raw_error"`
	PossibleReason   string `json:"possible_reason"`
}

// Update sets new data in structure.
func (debug *Debug) Update(status int, loc string, err string, reason string) {
	debug.StatusCode = status
	debug.Location = loc
	debug.RawError = err
	//update reason if status code shows client error
	if status == http.StatusBadRequest || status == http.StatusNotFound || status == http.StatusMethodNotAllowed {
		debug.PossibleReason = reason
	} else {
		debug.PossibleReason = "Unknown"
	}
}

// Print sends structure to client and console.
func (debug *Debug) Print(w http.ResponseWriter) {
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(debug.StatusCode)
	//send output to user and branch if an error occured
	err := json.NewEncoder(w).Encode(debug)
	if err != nil {
		fmt.Println("Error encoding JSON in Debug.Print()", err)
		return
	}
	//send output to console
	fmt.Printf(
		"%v {\n" +
		"    status_code:     %v,\n" +
		"    location:        %s,\n" +
		"    raw_error:       %s,\n" +
		"    possible_reason: %s \n" +
		"}\n", 
		time.Now().Format("2006-01-02 15:04:05"), debug.StatusCode, debug.Location, debug.RawError, debug.PossibleReason,
	)
}
