package diag

type DiagStatuses struct {
	Restcountries  int `json:"restcountries"`
	TicketMaster   int `json:"ticketmaster"`
	LocationIq     int `json:"locationiq"`
	Weatherapi     int `json:"weatherapi"`
	PublicHolidays int `json:"publicholidays"`

	RegisteredWebhooks int    `json:"registeredwebhooks"`
	Version            string `json:"version"`
	Uptime             int    `json:"uptime"`
}
