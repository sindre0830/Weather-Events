package countryData

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/api"
	"net/http"
	"os"
	"strings"
)

//Handler - Class function will be called and handle all requests and fetches
func (data Information) Handler(country string) (int, error, string) {

	//Check if document exists, if it does not exist run this code:
	if _, err := os.Stat("./data/countries.json"); os.IsNotExist(err) {

		//Fetch information:
		var input Information
		status, err := input.req("https://restcountries.eu/rest/v2/all")
		if err != nil {
			return status, err, ""
		}

		//Store it in a file:
		file, _ := json.MarshalIndent(input, "", " ")
		_ = ioutil.WriteFile("./data/countries.json", file, 0644)
	}

	//If function was called with a specific country
	if country != "" {
		var specificCountry Information
		status, err, countryCode := specificCountry.oneCountry(country)
		if err != nil {

			return status, err, ""
		}
		return http.StatusOK, nil, countryCode

	} else { //If function was NOT called with a specific country:
		var AllData Information
		status, err := AllData.allCountries()
		if err != nil {

			return status, err, ""
		}
		return http.StatusOK, nil, ""
	}
}

//req -Requests information from the api
func (data *Information) req(url string) (int, error) {
	output, status, jsonErr := api.RequestData(url)

	if jsonErr != nil {
		return status, jsonErr
	}

	jsonErr = json.Unmarshal(output, &data)
	if jsonErr != nil {
		return http.StatusInternalServerError, jsonErr
	}
	return http.StatusOK, jsonErr
}

//oneCountry -Gets information about one specific country from local storage
func (data *Information) oneCountry(countryName string) (int, error, string) {
	countryName = strings.Title(strings.ToLower(countryName))
	//Get data from document
	file, err := ioutil.ReadFile("./data/countries.json")
	if err != nil {
		return http.StatusInternalServerError, err, ""
	}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return http.StatusInternalServerError, err, ""
	}

	//Looping to look for the specific country
	for _, v := range *data {
		if v.Name == countryName {
			//When found, break the loop and return that one instance
			*data = Information{{v.Name, v.Alpha2Code, v.Capital}}
			return http.StatusOK, nil, v.Alpha2Code
		}
	}
	return http.StatusBadRequest, &MyError{"Country:" + countryName + " was not found in the database"}, ""
}

//allCountries -Gets information about all countries from local storage
func (data *Information) allCountries() (int, error) {
	//Get data from document
	file, err := ioutil.ReadFile("./data/countries.json")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

//ParseFile -From restStub (https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/tree/master/RESTstub)
func ParseFile(filename string) ([]byte, int, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return file, http.StatusInternalServerError, err
	}
	return file, http.StatusOK, err
}

//Error -Custom error function
func (e *MyError) Error() string {
	return fmt.Sprintf("Error: %s",
		e.What)
}
