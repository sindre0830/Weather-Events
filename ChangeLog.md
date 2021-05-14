# Changelog for weather-events

### Legend

[milestones.development.setup/formatting/fix]   **type of commit**:     Commit message

[0.0.1]     **setup**:          Commit about repository

[0.0.2]     **formatting**:     Commit about documentation

[0.0.3]     **fix**:            Commit about fixing bugs

[0.1.0]     **development**:    Commit about implimenting features

[1.0.0]     **milestone**:      Commit message about what has been implimented so far

### Logs

[0.0.0]     **setup**:          Initial commit

[0.0.1]     **setup**:          Initial commit of main and dictionary

[0.0.2]     **setup**:          Initial commit of dataHandling (from assignment 2)

[0.0.3]     **setup**:          Initial commit of error handler (from assignment 2)

[0.0.4]     **setup**:          Fixed ioutil.ReadAll, dataHandling.go name

[0.0.5]     **setup**:          Added endpoint description

[0.0.6]     **setup**:          Added database package (from assignment 2)

[0.1.0]     **development**:    Added database setup in init

[0.2.0]     **development**:    Added initial version of Yr structure

[0.3.0]     **development**:    Implemented get and req functors to Yr structure

[0.4.0]     **development**:    Added initial version of WeatherData structure

[0.5.0]     **development**:    Modified Yr structure to handle more context and modified get functionality to be local

[0.6.0]     **development**:    Modified WeatherData structure to use Camel Case naming style

[0.7.0]     **fix**:            Fixed Yr structure by nesting the Data structure

[0.8.0]     **development**:    Implemented get functor to WeatherData structure

[0.9.0]     **development**:    Added functionality where information about all countries (or filtered by one country) is fetched

[0.10.0]    **fix**:            Fixed Yr structure by adding type to Data field

[0.11.0]    **development**:    Implemented Handle functor to WeatherData structure

[0.12.0]    **development**:    Implemented method handler for WeatherData

[0.13.0]    **development**:    Added User-Agent with project information when requesting data (allows yr requests)

[0.14.0]    **development**:    HandlerCoords file with functions that handle the locationiq api datasource.

[0.15.0]    **fix**:            Restructured countrydata to be its own package + minor fixes related to errorhandling

[0.16.0]    **formatting**:     Small fix on endpoint in main.go

[0.17.0]    **development**:    HandlerCoords errors fixed up.

[0.18.0]    **development**:    Added standard data structure for database with container for different structures

[0.19.0]    **development**:    Added function to get all holidays of a country

[0.20.0]    **development**:    Added possiblity for custom id in database

[0.21.0]    **development**:    WeatherData Handler sends data to database

[0.22.0]    **development**:    Added database reader

[0.23.0]    **development**:    Added time validation on data stored in database

[0.24.0]    **development**:    Refactored HandlerCoords to check firestore for existing data.

[0.25.0]    **development**:    Fixed time validation when errors occur

[0.26.0]    **development**:    Added updated field in WeatherData

[0.27.0]    **development**:    Added ability to add and get holidays from the database

[0.28.0]    **development**:    Saving information from restCountries localy instead of firebase

[0.29.0]    **development**:    Added reading from and writing to local database in addition to firestore in HandlerCoords

[1.0.0]     **milestone**:      Implimented services to build our endpoints

[1.1.0]     **development**:    Added inital version of Weather structure

[1.2.0]     **development**:    Modified WeatherData to not be an endpoint

[1.3.0]     **development**:    Implemented get functor to Weather structure

[1.4.0]     **development**:    Added country names to LocationCoords struct, reimplemented handler as class method

[1.5.0]     **development**:    Reimplemented rest country handler as class method

[1.6.0]     **development**:    Reimplemented holidays handler as class method

[1.7.0]     **development**:    Changed how Handler in geocoords works, removed MethodHandler, changed Country to Address

[1.8.0]     **development**:    Changed errorhandling a little in restcountries and added some more comments

[1.9.0]     **development**:    Implemented handler functor and method handler for Weather structure

[1.9.1]     **fix**:            Fixed file reading for LocationCoords structure

[1.10.0]    **development**:    Added inital version of WeatherCompare structure

[1.10.1]    **formatting**:     Reformated Handler in geocoords to be cleaner

[1.11.0]    **development**:    Changed holidaysData so it checks if the year stored is the current

[1.12.0]    **development**:    Added fun package with decimal limiter

[1.13.0]    **development**:    Implemented get functor to WeatherCompare structure

[1.13.1]    **fix**:            Fixed database reading for LocationCoords structure

[1.14.0]    **development**:    Implemented handler functor and method handler for WeatherCompare structure

[1.15.0]    **development**:    Added information about endpoints

[1.16.0]    **development**:    Removed filtering suggestion from Weather and WeatherCompare

[2.0.0]     **milestone**:      Implemented all main endpoints

[2.1.0]     **development**:    Moved information about used REST services

[2.2.0]     **development**:    Added more comments on restCountry file, edited readme

[2.2.1]     **reformatting**:   Moved fields around in Weather and WeatherCompare structure and added comments to structure functors

[2.3.0]     **development**:    Modified information about endpoints

[2.4.0]     **development**:    Implemented method handler and some basic functionality of WeatherHoliday structure

[2.5.0]     **development**:    Implemented function to get a country's alpha code

[2.6.0]     **development**:    Implemented function to add weatherHoliday webhooks to the database

[2.7.0]     **development**:    Renamed weather fields Now and Today to Instant and Predicted

[2.8.0]     **development**:    Simplified data reading in WeatherData

[2.9.0]     **development**:    Modified WeatherData to store data for all available days 

[2.10.0]    **development**:    Initial modification of Weather and WeatherCompare to handle inputted dates

[2.11.0]    **development**:    Added handling of inputted date for Weather structure

[2.12.0]    **development**:    Added handling of inputted date for WeatherCompare structure

[2.13.0]    **development**:    Added function to delete weatherHoliday webhook

[2.14.0]    **development**:    Changed adding to database function to return the documents ID

[2.15.0]    **development**:    Added getWeatherURL in dictionary

[2.16.0]    **development**:    Added getWeatherCompareURL in dictionary

[2.17.0]    **development**:    Added ticketmaster endpoint

[2.18.0]    **development**:    Moved code from holiday webhook to holiday data

[2.18.1]    **fix**:            Fixed an error when formatting holiday name wrong

[2.19.0]    **development**:    Added WeatherEvent POST method and a feedback pipeline to handle webhook feedback

[2.20.0]    **development**:    Added passing of data from weatherHoliday to weatherEvent

[2.21.0]    **development**:    Added WeatherEvent GET method and modified database.Get() to return map of interfaces

[2.21.1]    **fix**:            Fixed weatherHook to work with new database.Get()

[2.21.2]    **development**:    Slight cleanup of db.Delete, more work on weatherHook

[2.21.3]    **fix**:            Fixed WeatherHoliday to work with new database.Get()

[2.22.0]    **development**:    Finished WeatherHoliday POST function 

[2.23.0]    **development**:    Added diag endpoint and restructured eventData

[2.24.0]    **development**:    Added WeatherEvent DELETE method

[2.25.0]    **development**:    Added WeatherEvent callLoop after a POST request

[2.26.0]    **development**:    Added function to get one or all WeatherHoliday webhooks

[2.26.1]    **fix**:            Solved issue with Postman seeing POST requests as GET

[2.26.2]    **fix**:            Fixed WeatherHoliday Get so the output is JSON

[2.26.3]    **fix**:            Fixed HandlerCoord's Handler to check if "Time" key exists

[2.27.0]    **development**:    Modified WeatherEvent callLoop to terminate when it is deleted

[2.28.0]    **development**:    Added StartTrigger to weatherHook for initiating webhooks from DB on program start

[2.28.1]    **fix**:            Fixed WeatherData's Handler to check if "Time" key exists

[2.29.0]    **development**:    Added function to check if Date field is date or holiday

[2.29.1]    **fix**:            Fixed formatting of date

[2.30.0]    **development**:    Changed Ticketmaster endpoint to be able to be called from another function, started moving structs to separete file

[2.31.0]    **development**:    Moved weatherHook, changed trigger to callUrl.

[2.32.0]    **development**:    Added date validation in weatherEvent.POST()

[2.33.0]    **development**:    Restructured so all structs have their own file in that same package

[2.34.0]    **development**:    Added date validation in weatherEvent.callLoop()

[2.34.1]    **formatting**:     Reformatted structures in WeatherData package

[2.34.2]    **formatting**:     Reformatted dict package

[2.34.3]    **formatting**:     Reformatted weatherData package

[2.34.4]    **formatting**:     Reformatted weather package

[2.34.5]    **fix**:            Fixed CheckDate to subtract hours instead of minutes

[2.34.6]    **formatting**:     Reformatted weatherCompare package

[2.34.7]    **fix**:            Fixed error handling in weatherEvent.checkIfHoliday()

[2.34.8]    **fix**:            Removed weatherHoliday package

[2.34.9]    **development**:    Separated sleep function in fun, added comments to weatherHook

[2.34.10]   **formatting**:     Added commentary on some files

[2.35.0]    **development**:    Added initialization of WeatherEvent hooks on startup

[2.35.1]    **formatting**:     Reformatted weatherEvent package

[2.36.0]    **development**:    Added mutex lock to weatherHook callUrl, big update to readme

[2.37.0]    **development**:    Fixed up weatherHook.callLoop function

[2.38.0]    **development**:    Added different types of WeatherEvent webhooks

[2.39.0]    **development**:    Added mutex lock to WeatherEvent

[2.40.0]    **development**:    Started adding unit testing

[2.41.0]    **development**:    Restructured webhook count in diag to reduce operations

[3.0.0]     **milestone**:      Implemented all webhooks

[3.1.0]     **development**:    Added database cleaner which removes some collections every 12 hours

[3.2.0]     **development**:    Renamed db package to storage and moved file operations there

[3.2.1]     **setup**:          Moved GeoCoord file to data/ and added a gitkeep

[3.3.0]     **development**:    Restructured Weather and WeatherCompare structures

[3.4.0]     **development**:    Renamed weather to weatherDetails and compare to weatherCompare (packages)

[3.5.0]     **development**:    Restructured Weather hook to be more similar to WeatherEvent, added a common mutex state to handle multiple webhooks

[3.6.0]     **development**:    Added more unit tests