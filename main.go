package main

import (
	"log"
	"main/api/holidays"
	"main/db"
	"main/dict"
	"net/http"
	"os"
)

// init runs once at startup.
func init() {
	//setup connection with firebase and branch if an error occured
	err := db.DB.Setup()
	if err != nil {
		defer db.DB.Client.Close()
		log.Fatalln(err)
	}
}

// Main program.
func main() {
	//get port and branch if there isn't a port and set it to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	//set URL with port
	dict.URL = dict.URL + ":" + port


	//handle weather data
	//http.HandleFunc("/weather-rest/v1/weather/data/", weatherData.MethodHandler)
	//Get all countries endpoint:
	//http.HandleFunc("/weather-rest/v1/restCountries/", countryData.HandleRestCountry)
	// Get a country's holidays
	http.HandleFunc("/weather-rest/v1/holidays/", holidays.GetCountryHolidays)

	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
