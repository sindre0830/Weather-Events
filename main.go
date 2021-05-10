package main

import (
	"log"
	"main/api/eventData"
	"main/api/notification/weatherEvent"
	"main/api/weather"
	compare "main/api/weatherCompare"
	"main/api/weatherHoliday"
	"main/api/weatherHook"
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
	http.HandleFunc(dict.WEATHER_PATH, weather.MethodHandler)
	//handle weather comparison data
	http.HandleFunc(dict.WEATHERCOMPARE_PATH, compare.MethodHandler)
	//handle weather event
	http.HandleFunc(dict.WEATHEREVENT_PATH, weatherEvent.MethodHandler)
	//handle event data
	http.HandleFunc(dict.EVENT_PATH, eventData.MethodHandler)
	//handle holiday webhook
	http.HandleFunc(dict.HOLIDAY_PATH, weatherHoliday.MethodHandler)
	//handle weather webhook
	http.HandleFunc(dict.WEATHERHOOK_PATH, weatherHook.MethodHandler)
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
