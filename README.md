# [Project | weather-events](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Project-Description)

### Info
- Authors: 
    - Sindre Eiklid (sindreik@stud.ntnu.no)
- Root path:
    - Main:     localhost:8080/weather-rest/v1
    - Client:   *TBA*
- We have used these REST web services to build our service:
    - *TBA*
- You need to be connected to NTNU network with a VPN to run the program. If you want to run it locally, you will have to change the URL variable in the 'dict' package to ```http://localhost```.
- Client Repo: *TBA*

### Usage

1. Weather

    - Input:
        ```
        Method: GET
        Path: .../weather/location/{:location}{?fields=filter1;filter2;...}
        ```
        - **{:location}** *TBA*
        - **{?fields=filter1;filter2;...}** *TBA*

    - Output:
        ```go
        type Weather struct {
            Location             string  `json:"location"`
            Longitude            float64 `json:"longitude"`
            Latiude              float64 `json:"latiude"`
            Altitude             float64 `json:"altitude"`
            Country              string  `json:"country"`
            Updated              string  `json:"updated"`
            Data struct {
                Now struct {
                    Air_temperature         float64 `json:"air_temperature"`
                    Cloud_area_fraction     float64 `json:"cloud_area_fraction"`
                    Dew_point_temperature   float64 `json:"dew_point_temperature"`
                    Relative_humidity       float64 `json:"relative_humidity"`
                    Wind_from_direction     float64 `json:"wind_from_direction"`
                    Wind_speed              float64 `json:"wind_speed"`
                    Wind_speed_of_gust      float64 `json:"wind_speed_of_gust"`
                } `json:"now"`
                Today struct {
                    Summary                         string  `json:"summary"`
                    Confidence                      string  `json:"confidence"`
                    air_temperature_max             float64 `json:"air_temperature_max"`
                    air_temperature_min             float64 `json:"air_temperature_min"`
                    precipitation_amount            float64 `json:"precipitation_amount"`
                    precipitation_amount_max        float64 `json:"precipitation_amount_max"`
                    precipitation_amount_min        float64 `json:"precipitation_amount_min"`
                    Probability_of_precipitation    float64 `json:"probability_of_precipitation"`
                } `json:"today"`
            } `json:"data"`
        }
        ```

    - Example:
        - Input: 
            *TBA*
        - Output:
            *TBA*

2. Compare

    - Input:
        ```
        Method: GET
        Path: .../weather/compare/{:location}/{:compare=location1;location2;...}{?fields=filter1;filter2;...}
        ```
        - **{:location}** *TBA*
        - **{?compare=location1;location2;...}** *TBA*
        - **{?fields=filter1;filter2;...}** *TBA*

    - Output:
        ```go
        type CompareWeather struct {
            Updated              string  `json:"updated"`
            Main_location        string  `json:"main_location"`
            data []struct {
                Location             string  `json:"location"`
                Longitude            float64 `json:"longitude"`
                Latiude              float64 `json:"latiude"`
                Altitude             float64 `json:"altitude"`
                Country              string  `json:"country"`
                Now struct {
                    Air_temperature         float64 `json:"air_temperature"`
                    Cloud_area_fraction     float64 `json:"cloud_area_fraction"`
                    Dew_point_temperature   float64 `json:"dew_point_temperature"`
                    Relative_humidity       float64 `json:"relative_humidity"`
                    Wind_speed              float64 `json:"wind_speed"`
                    Wind_speed_of_gust      float64 `json:"wind_speed_of_gust"`
                } `json:"now"`
                Today struct {
                    air_temperature_max             float64 `json:"air_temperature_max"`
                    air_temperature_min             float64 `json:"air_temperature_min"`
                    precipitation_amount            float64 `json:"precipitation_amount"`
                    precipitation_amount_max        float64 `json:"precipitation_amount_max"`
                    precipitation_amount_min        float64 `json:"precipitation_amount_min"`
                    Probability_of_precipitation    float64 `json:"probability_of_precipitation"`
                } `json:"today"`
            } `json:"data"`
        }
        ```

    - Example:
        - Input: 
            *TBA*
        - Output:
            *TBA*

3. Holidays

    - Input:
        ```
        Method: GET
        Path: .../events/holidays/{:location}{?holiday=holiday}
        ```
        - **{:location}** *TBA*
        - **{?holiday=holiday}** *TBA*

    - Output:
        ```go
        type Holidays struct {
            Updated              string             `json:"updated"`
            Location             string             `json:"location"`
            Holiday              map[string]string  `json:"holiday"`
        }
        ```

    - Example:
        - Input: 
            *TBA*
        - Output:
            *TBA*

## Notes

#### Design Decisions

*TBA*

#### Structure

├──api
│   ├── countryData
│   │   └── restCountries.go
│   ├── geoCoordsData
│   │   └── HandlerCoords.go
│   ├── holidaysData
│   │   └── holidays.go
│   ├── weather
│   │   ├── methodHandler.go
│   │   └── weather.go
│   ├── weatherCompare
│   │   ├── methodHandler.go
│   │   └── weatherCompare.go
│   ├── weatherData
│   │   ├── weatherData.go
│   │   └── yr.go
│   └── dataHandling
├── db
│   └── database.go
├── debug
│   └── errorHandler.go
├── dict
│   └── dictionary.go
├── fun
│   └── math.go
├── ChangeLog.md
├── go.mod
├── go.sum
├── main.go
├── README.md
└── fun

#### Error Handling

When an error is found, we add info to a debug struct.
Debugging information is then sent to user as a json object, and printed to console.

```go
	StatusCode 		 int    `json:"status_code"`                // The REST code for the error
	Location   		 string `json:"location"`                   // Where in the program did the error occur
	RawError   		 string `json:"raw_error"`                  // The raw error data
	PossibleReason   string `json:"possible_reason"`            // Potential reasons for the error occurring (e.g. misspelled endpoint, etc)
```

#### Testing

*TBA*

##### Usage
For Visual Studio Code with Golang extension:
1. Open testing file in the IDE
2. Click the ```run test``` label for any function that you want to test
