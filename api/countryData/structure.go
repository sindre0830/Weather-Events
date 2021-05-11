package countryData

//Information -Struct containing all necessary information from restcountries
type Information []struct {
	Name       string `json:"name"`
	Alpha2Code string `json:"alpha2Code"`
	Capital    string `json:"capital"`
}

//MyError -Struct containing custom error message
type MyError struct {
	What string
}
