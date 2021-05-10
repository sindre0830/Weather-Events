package notification

import (
	"encoding/json"
	"net/http"
)

// Feedback structure stores information about successful method request.
//
// Functionality: update, print
type Feedback struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ID		   string `json:"id"`
}

// Update sets new data in structure.
func (feedback *Feedback) Update(status int, message string, id string) {
	feedback.StatusCode = status
	feedback.Message = message
	feedback.ID = id
}

// Print sends structure to client.
func (feedback *Feedback) Print(w http.ResponseWriter) error {
	//update header to JSON and set HTTP code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(feedback.StatusCode)
	//send output to user and branch if an error occured
	err := json.NewEncoder(w).Encode(feedback)
	if err != nil {
		return err
	}
	return nil
}
