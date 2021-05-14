package notification

import (
	"main/api"
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
	//send output to user and return if an error occurred
	return api.SendData(w, feedback, feedback.StatusCode)
}
