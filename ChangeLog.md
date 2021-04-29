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

[1.9.2]     **formatting**:     Reformated Handler in geocoords to be cleaner
