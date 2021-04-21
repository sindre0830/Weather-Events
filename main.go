package main

import (
	"log"
	"main/dict"
	"net/http"
	"os"
)

// init runs once at startup.
func init() {

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
	//ends program if it can't open port
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
