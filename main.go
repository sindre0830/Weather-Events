package main

import (
	"log"
	"main/api/eventData"
	"main/api/weather"
	compare "main/api/weatherCompare"
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
	http.HandleFunc("/weather-rest/v1/weather/location/", weather.MethodHandler)
	//handle weather comparison data
	http.HandleFunc("/weather-rest/v1/weather/compare/", compare.MethodHandler)
	//handle event data
	http.HandleFunc("/weather-rest/v1/weather/event/", eventData.MethodHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
