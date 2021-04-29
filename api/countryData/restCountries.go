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

//Struct containing all information from restcountries
type Information []struct {
	Name       string `json:"name"`
	Alpha2Code string `json:"alpha2Code"`
	Capital    string `json:"capital"`
}

//Struct containing error message
type MyError struct {
	What string
}

// Class function will be called and handle all requests and fetches
func (data *Information) Handler(country string) (int, error) {
	//Check if document exists
	if _, err := os.Stat("./data/countries.json"); os.IsNotExist(err) { //if it does not exist:
		//Fetch information
		var input Information
		status, err := input.req("https://restcountries.eu/rest/v2/all")
		if err != nil {
			return status, err
		}

		//Store it in a file:
		file, _ := json.MarshalIndent(input, "", " ")

		_ = ioutil.WriteFile("./data/countries.json", file, 0644)
	}

	//If you want a specific country
	if country != "" {
		var specificCountry Information
		status, err := specificCountry.oneCountry(country)
		if err != nil {

			return status, err
		}
		return http.StatusOK, nil

	} else { //if you dont require a specific country:
		var AllData Information
		status, err := AllData.allCountries()
		if err != nil {

			return status, err
		}
		return http.StatusOK, nil
	}
}

//req requests information
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

//oneCountry gets information about one specific country
func (data *Information) oneCountry(countryName string) (int, error) {
	countryName = strings.Title(strings.ToLower(countryName))
	//Get data from document
	file, err := ioutil.ReadFile("./data/countries.json")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, v := range *data {
		if v.Name == countryName {

			//only return one data instance
			*data = Information{{v.Name, v.Alpha2Code, v.Capital}}

			return http.StatusOK, nil
		}

	}

	return http.StatusBadRequest, &MyError{"Country:" + countryName + " was not found in the database"}
}

//allCountries gets information about all countries
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

//From restStub (https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/tree/master/RESTstub)
func ParseFile(filename string) ([]byte, int, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return file, http.StatusInternalServerError, err
	}
	return file, http.StatusOK, err
}

//Custom error function
func (e *MyError) Error() string {
	return fmt.Sprintf("Error: %s",
		e.What)
}
