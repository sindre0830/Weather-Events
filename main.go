package main

import (
	"log"
	"main/api/diag"
	"main/api/notification/weather"
	"main/api/notification/weatherEvent"
	"main/api/weatherCompare"
	"main/api/weatherDetails"
	"main/dict"
	"main/storage"
	"net/http"
	"os"
	"time"
)

// init runs once at startup.
func init() {
	//start timer
	diag.StartTime = time.Now()
	//setup connection with firebase and branch if an error occured
	err := storage.Firebase.Setup()
	if err != nil {
		defer storage.Firebase.Client.Close()
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
	dict.MAIN_URL = dict.MAIN_URL + ":" + port
	//start webhooks
	weatherEvent.InitHooks()
	weather.InitHooks()
	//handle weather data
	http.HandleFunc(dict.WEATHERDETAILS_PATH, weatherDetails.MethodHandler)
	//handle weather comparison data
	http.HandleFunc(dict.WEATHERCOMPARE_PATH, weatherCompare.MethodHandler)
	//handle weather event
	http.HandleFunc(dict.WEATHEREVENT_HOOK_PATH, weatherEvent.MethodHandler)
	//Diag endpoint
	http.HandleFunc(dict.DIAG_PATH, diag.MethodHandler) //NB Note that the count of webhooks counts the collections, therefore they need to be added manually and as such not all webhooks are counted as of yet
	//handle weather webhook
	http.HandleFunc(dict.WEATHER_HOOK_PATH, weather.MethodHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
