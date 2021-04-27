package main

import (
	"log"
	"main/api"
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

	//Get all countries endpoint:
	http.HandleFunc("/weather/v1/restCountries/", api.HandleRestCountry)

	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
