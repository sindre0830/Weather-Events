package main

import (
	"log"
	"main/api/diag"
	"main/api/notification/weatherEvent"
	"main/api/notification/weatherHook"
	"main/api/weather"
	compare "main/api/weatherCompare"
	"main/db"
	"main/dict"
	"net/http"
	"os"
	"time"
)

// init runs once at startup.
func init() {
	//start timer
	diag.StartTime = time.Now()
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
	weatherHook.StartCall(&db.DB) // Can't do this in database.go - cycling imports
	//set URL with port
	dict.MAIN_URL = dict.MAIN_URL + ":" + port
	//handle weather data
	http.HandleFunc(dict.WEATHER_PATH, weather.MethodHandler)
	//handle weather comparison data
	http.HandleFunc(dict.WEATHERCOMPARE_PATH, compare.MethodHandler)
	//handle weather event
	http.HandleFunc(dict.WEATHEREVENT_PATH, weatherEvent.MethodHandler)
	//Diag endpoint
	http.HandleFunc(dict.DIAG_PATH, diag.MethodHandler) //NB Note that the count of webhooks counts the collections, therefore they need to be added manually and as such not all webhooks are counted as of yet
	//handle weather webhook
	http.HandleFunc(dict.WEATHERHOOK_PATH, weatherHook.MethodHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
