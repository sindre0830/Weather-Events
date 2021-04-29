package countryData

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/api"
	"main/debug"
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

type MyError struct {
	What string
}

func (data *Information) Handler(country string) (int, error) {
	//Check if document exists

	if _, err := os.Stat("./data/countries.json"); os.IsNotExist(err) { //if it does not exist:
		//Fetch information
		var input Information
		status, err := input.req("https://restcountries.eu/rest/v2/all")
		if err != nil {
			debug.ErrorMessage.Update(
				status,
				"CountryData.Handler() -> Getting country data",
				err.Error(),
				"Unknown",
			)
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
			debug.ErrorMessage.Update(
				status,
				"CountryData.Handler() ->  Getting specific country data",
				err.Error(),
				"Unknown",
			)
			return status, err
		}
		return status, err
	} else { //if you dont require a specific country:
		var AllData Information
		status, err := AllData.allCountries()
		if err != nil {
			debug.ErrorMessage.Update(
				status,
				"CountryData.Handler() ->  Getting specific country data",
				err.Error(),
				"Unknown",
			)
			return status, err
		}
		return status, err
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

//req requests information
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

//req requests information
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

//From restStub
func ParseFile(filename string) []byte {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	return file
}

//Custom error function
func (e *MyError) Error() string {
	return fmt.Sprintf("Error: %s",
		e.What)
}
