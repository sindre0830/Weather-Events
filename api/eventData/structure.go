package eventData

//FirebaseStore -simple struct containing ticketmaster information
type FirebaseStore struct {
	Localdate string `json:"localDate"`
	Name      string `json:"name"`
}

type EventInformation struct {
	Dates struct {
		Start struct {
			Localdate string `json:"localDate"`
		} `json:"start"`
	} `json:"dates"`

	Embedded struct {
		Venues []struct {
			City struct {
				Name string `json:"name"`
			} `json:"city"`
		} `json:"venues"`
	} `json:"_embedded"`
}
