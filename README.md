# [Project | weather-events](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Project-Description)

### Info
- Authors: 
    - Sindre Eiklid (sindreik@stud.ntnu.no)
    - Rickard Loland (rickarl@stud.ntnu.no)
    - Susanne Skjold Edvardsen (susanse@stud.ntnu.no)
    - Maren Skårestuen Grindal (marensg@stud.ntnu.no)
- Root path:
    - Main:     localhost:8080/weather-rest/v1
    - Client:   *TBA*
- We have used these REST web services to build our service:
    - Weather information:  https://api.met.no/weatherapi/locationforecast/2.0/complete
    - Country information:  https://restcountries.eu/rest/v2/all
    - Holiday information:  https://date.nager.at/api/v2/PublicHolidays/
    - Location information: https://us1.locationiq.com/v1/
- You need to be connected to NTNU network with a VPN to run the program. If you want to run it locally, you will have to change the URL variable in the 'dict' package to ```http://localhost```.
- Client Repo: *TBA*

### About:

The idea of this project is to utilize the weather data at yr's API and match it with event-based APIs (concerts, games, whatever) to let users find the weather for the event time + location. There are two endpoints, and one webhook. The first endpoint gives a basic weather report for the location, and the second compares a base location with other locations. The webhook lets a user pass in a location and a holiday, and the service will return the weather report for that date. The webhook gives the option to register for a future event and be notefied when the weather report for that event changes, or is updated for the first time.

In addition to these services another webhook may be implemented if there is time for it:

This webhook will allow a user to register an event id, like a concert. It will return the location and date of the event, additionally the weather if it allows for it, although weather is only availale 9 days ahead in time so a notification will be sent when its available. 

#### Progress

So far we have implemented most of the functionality. All endpoints except webhooks are done. It has gone smooth so far, working on incrementing each endpoint when natural and building on what we have so far. We are designing it to be easy to rewrite and repurpose. We are implementing helping functions and packages where it's fitting, and class methods for structs are quite numerous throughout the project. Working in a group has worked out well so far. We have had regular meetings and a structured plan which made it easy to actually get things done. This worked especially well while working on retrieving data from our service endpoints, as everyone could work siultaneously without issue. For some more difficult work, we brainstormed solutions together during meetings while one person implemented and pushed the code. We did have occasional bottlenecks where some of us had to wait for someone else to finish, but there was always refactoring, readme improvements and other things to fix up. None of these lasted very long, so they did not present a challenge for the project as a whole.

#### Experiences

The only problem we encounterd so far during development was almost exceeding the free firestore operation quota of 50K reads. We planned to store all of our data in firestore, but realized that each firestore query would read through every ID. This resoluted in each query reading more than 200 times. We ended up using more than 16% of our free quota in 2 days, after fixing this we used less than 0.5% each day.

![Firestore operations reaching 5.8K a day](images/firestore.png)

### Usage

1. Weather

    - Input:
        ```
        Method: GET
        Path: .../weather/location/{:location}
        ```

    - Output:
        ```go
        type Weather struct {
            Longitude float64 `json:"longitude"`
            Latitude  float64 `json:"latitude"`
            Location  string  `json:"location"`
            Updated   string  `json:"updated"`
            Data      struct {
                Now struct {
                    AirTemperature      float64 `json:"air_temperature"`
                    CloudAreaFraction   float64 `json:"cloud_area_fraction"`
                    DewPointTemperature float64 `json:"dew_point_temperature"`
                    RelativeHumidity    float64 `json:"relative_humidity"`
                    WindFromDirection   float64 `json:"wind_from_direction"`
                    WindSpeed           float64 `json:"wind_speed"`
                    WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
                    PrecipitationAmount float64 `json:"precipitation_amount"`
                } `json:"now"`
                Today struct {
                    Summary                    string  `json:"summary"`
                    Confidence                 string  `json:"confidence"`
                    AirTemperatureMax          float64 `json:"air_temperature_max"`
                    AirTemperatureMin          float64 `json:"air_temperature_min"`
                    PrecipitationAmount        float64 `json:"precipitation_amount"`
                    PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
                    PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
                    ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
                } `json:"today"`
            } `json:"data"`
        }
        ```

    - Example:
        - Input: 
            ```
            Method: GET
            Path: localhost:8080/weather-rest/v1/weather/location/oslo
            ```
        - Output:
            ```json
            {
                "longitude": 10.74,
                "latitude": 59.91,
                "location": "Oslo, 0026, Norway",
                "updated": "29 Apr 21 11:20 CEST",
                "data": {
                    "now": {
                        "air_temperature": 9.2,
                        "cloud_area_fraction": 7,
                        "dew_point_temperature": -4.4,
                        "relative_humidity": 38.6,
                        "wind_from_direction": 57.8,
                        "wind_speed": 5.1,
                        "wind_speed_of_gust": 8.8,
                        "precipitation_amount": 0
                    },
                    "today": {
                        "summary": "fair_day",
                        "confidence": "certain",
                        "air_temperature_max": 12.2,
                        "air_temperature_min": 10.1,
                        "precipitation_amount": 0,
                        "precipitation_amount_max": 0,
                        "precipitation_amount_min": 0,
                        "probability_of_precipitation": 0
                    }
                }
            }
            ```

2. Compare

    - Input:
        ```
        Method: GET
        Path: .../weather/compare/{:location}/{:location1;location2;...}
        ```

    - Output:
        ```go
        type WeatherCompare struct {
            Longitude float64 `json:"longitude"`
            Latitude  float64 `json:"latitude"`
            Location  string  `json:"location"`
            Updated   string  `json:"updated"`
            Data         struct {
                Longitude float64 `json:"longitude"`
                Latitude  float64 `json:"latitude"`
                Location  string  `json:"location"`
                Updated   string  `json:"updated"`
                Now       struct {
                    AirTemperature      float64 `json:"air_temperature"`
                    CloudAreaFraction   float64 `json:"cloud_area_fraction"`
                    DewPointTemperature float64 `json:"dew_point_temperature"`
                    RelativeHumidity    float64 `json:"relative_humidity"`
                    WindSpeed           float64 `json:"wind_speed"`
                    WindSpeedOfGust     float64 `json:"wind_speed_of_gust"`
                    PrecipitationAmount float64 `json:"precipitation_amount"`
                } `json:"now"`
                Today struct {
                    AirTemperatureMax          float64 `json:"air_temperature_max"`
                    AirTemperatureMin          float64 `json:"air_temperature_min"`
                    PrecipitationAmount        float64 `json:"precipitation_amount"`
                    PrecipitationAmountMax     float64 `json:"precipitation_amount_max"`
                    PrecipitationAmountMin     float64 `json:"precipitation_amount_min"`
                    ProbabilityOfPrecipitation float64 `json:"probability_of_precipitation"`
                } `json:"today"`
            } `json:"data"`
        }
        ```

    - Example:
        - Input: 
            ```
            Method: GET
            Path: localhost:8080/weather-rest/v1/weather/compare/oslo/bergen;stavanger
            ```
        - Output:
            ```json
            {
                "longitude": 10.74,
                "latitude": 59.91,
                "location": "Oslo, 0026, Norway",
                "updated": "29 Apr 21 11:20 CEST",
                "data": [
                    {
                        "longitude": 5.33,
                        "latitude": 60.39,
                        "location": "Bergen, Vestland, Norway",
                        "updated": "29 Apr 21 12:26 CEST",
                        "now": {
                            "air_temperature": 0.1,
                            "cloud_area_fraction": 11.3,
                            "dew_point_temperature": 1.3,
                            "relative_humidity": 3.8,
                            "wind_speed": -2.4,
                            "wind_speed_of_gust": -2.7,
                            "precipitation_amount": 0
                        },
                        "today": {
                            "air_temperature_max": -1.6,
                            "air_temperature_min": -0.1,
                            "precipitation_amount": 0,
                            "precipitation_amount_max": 0,
                            "precipitation_amount_min": 0,
                            "probability_of_precipitation": 0
                        }
                    },
                    {
                        "longitude": 5.71,
                        "latitude": 59.1,
                        "location": "Stavanger, Rogaland, Norway",
                        "updated": "29 Apr 21 13:09 CEST",
                        "now": {
                            "air_temperature": 1.9,
                            "cloud_area_fraction": 31.5,
                            "dew_point_temperature": 5.2,
                            "relative_humidity": 10.3,
                            "wind_speed": -3.3,
                            "wind_speed_of_gust": -0.6,
                            "precipitation_amount": 0
                        },
                        "today": {
                            "air_temperature_max": -0.1,
                            "air_temperature_min": 1.2,
                            "precipitation_amount": 0,
                            "precipitation_amount_max": 0,
                            "precipitation_amount_min": 0,
                            "probability_of_precipitation": 0
                        }
                    }
                ]
            }
            ```

## Notes

#### Design Decisions

##### Technologies used

The technologies we are going to use are Firestore, OpenStack and Docker. We are using firestore for caching. The weather data is stored for 6 hours. Whether or not the geocoords are stored depends on the importance of the selected location. If it has a low importance, it is stored for 3 hours. If the importance is high, it is saved in a file. The data about holidays are stored until the year change. 

#### Structure

```
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
```

#### Error Handling

When an error is found, we add info to a debug struct with an Update struct function.
Debugging information is then sent to user as a json object, and printed to console with a Print function.

```go
type Debug struct {
	StatusCode       int    `json:"status_code"`     // The REST code for the error
	Location         string `json:"location"`        // Where in the program did the error occur
	RawError         string `json:"raw_error"`       // The raw error data
	PossibleReason   string `json:"possible_reason"` // Potential reasons for the error occurring (e.g. misspelled endpoint, etc)
}
```

#### Testing

*TBA*

##### Usage
For Visual Studio Code with Golang extension:
1. Open testing file in the IDE
2. Click the ```run test``` label for any function that you want to test
