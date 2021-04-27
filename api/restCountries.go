//link https://restcountries.eu/rest/v2/all

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

//Struct containing all information from restcountries
type Data []struct {
	Name           string    `json:"name"`
	Topleveldomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	Callingcodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	Altspellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	Nativename     string    `json:"nativeName"`
	Numericcode    string    `json:"numericCode"`
	Currencies     []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Iso6391    string `json:"iso639_1"`
		Iso6392    string `json:"iso639_2"`
		Name       string `json:"name"`
		Nativename string `json:"nativeName"`
	} `json:"languages"`
	Translations struct {
		De string `json:"de"`
		Es string `json:"es"`
		Fr string `json:"fr"`
		Ja string `json:"ja"`
		It string `json:"it"`
		Br string `json:"br"`
		Pt string `json:"pt"`
		Nl string `json:"nl"`
		Hr string `json:"hr"`
		Fa string `json:"fa"`
	} `json:"translations"`
	Flag          string `json:"flag"`
	Regionalblocs []struct {
		Acronym       string        `json:"acronym"`
		Name          string        `json:"name"`
		Otheracronyms []interface{} `json:"otherAcronyms"`
		Othernames    []interface{} `json:"otherNames"`
	} `json:"regionalBlocs"`
	Cioc string `json:"cioc"`
}

func HandleRestCountry(w http.ResponseWriter, r *http.Request) {
	//url parsing
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		status := http.StatusBadRequest
		http.Error(w, "Format was wrong, expecting /weather/v1/restCountries?specificCountry={:countryname}", status) //Input didnt match with was expected
		return
	}

	query := r.URL.Query()
	country := query.Get("specificCountry")
	var input Data
	if country == "" {
		status, err := input.req("https://restcountries.eu/rest/v2/all")
		if err != nil {
			log.Println("Error in fetching data from API", err)
			http.Error(w, "Error: ", status)
			return
		}
	} else {
		status, err := input.req("https://restcountries.eu/rest/v2/name/" + country + "?fullText=true")
		if err != nil {
			log.Println("Error in fetching data from API", err)
			http.Error(w, "Error: ", status)
			return
		}
	}

	//Formats the printouts
	w.Header().Set("Content-Type", "application/json")
	//Outputs results
	err := json.NewEncoder(w).Encode(input)
	if err != nil {
		log.Println("ERROR encoding JSON", err) //If error, send log to error logger
	}
}

//req requests information
func (data *Data) req(url string) (int, error) {
	output, status, jsonErr := RequestData(url)

	if jsonErr != nil {
		return status, jsonErr
	}

	jsonErr = json.Unmarshal(output, &data)
	if jsonErr != nil {
		return http.StatusInternalServerError, jsonErr
	}
	return http.StatusOK, jsonErr
}
