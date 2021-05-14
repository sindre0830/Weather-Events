package weather

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/api/diag"
	"main/api/weatherDetails"
	"main/dict"
	"main/storage"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

// InitHooks initilizes all weather hooks from the database.
func InitHooks() {
	//get all webhooks and branch if an error occured
	arrWeather, err := storage.Firebase.GetAll(dict.WEATHER_COLLECTION)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when initializing Weather webhooks.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		return
	}
	//call each hook
	for _, data := range arrWeather {
		var weather Weather
		weather.readData(data["Container"].(interface{}))
		go weather.callHook()
	}
	//print message with amount of webhooks initilizied
	fmt.Printf(
		"%v {\n\tSuccesfully initialized Weather webhooks.\n\tAmount: %v\n}\n",
		time.Now().Format("2006-01-02 15:04:05"), strconv.Itoa(len(arrWeather)),
	)
	diag.HookAmount += len(arrWeather)
}

// callHook calls webhook.
func (weather *Weather) callHook() {
	//check if webhook still exist in database
	_, exist := storage.Firebase.Get(dict.WEATHER_COLLECTION, weather.ID)
	if !exist {
		return
	}
	//check if program should sleep on timeout value
	nextTime := time.Now().Truncate(time.Second)
	nextTime = nextTime.Add(time.Duration(weather.Timeout) * time.Second)
	time.Sleep(time.Until(nextTime))
	dict.MutexState.Lock()
	//create new GET request and branch if an error occurred
	var weatherDetails weatherDetails.WeatherDetails
	req, err := http.NewRequest(http.MethodGet, dict.GetWeatherDetailsURL(weather.Location, ""), nil)
	if err != nil {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
			"Error when creating HTTP request to Weather.Handler().", 
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return 
	}
	//call the policy handler and branch if the status code is not OK
	//this stops timed out request being sent to the webhook
	recorder := httptest.NewRecorder()
	weatherDetails.Handler(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, recorder.Result().StatusCode, "None",
			"Error when creating HTTP request to Weather.Handler().", 
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//convert from structure to bytes and branch if an error occurred
	output, err := json.Marshal(weatherDetails)
	if err != nil {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
			"Error when parsing Weather structure.", 
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//create new POST request and branch if an error occurred
	req, err = http.NewRequest(http.MethodPost, weather.URL, bytes.NewBuffer(output))
	if err != nil {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
			"Error when creating new POST request.", 
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//hash structure and branch if an error occurred
	mac := hmac.New(sha256.New, dict.Secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
			"Error when hashing content before POST request.", 
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//convert hashed structure to string and add to header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Signature", hex.EncodeToString(mac.Sum(nil)))
	//send request to client and branch if an error occured
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
			"Error when sending HTTP content to webhook. Putting webhook to sleep for 6 hours...", 
		)
		dict.MutexState.Unlock()
		time.Sleep(time.Duration(6) * time.Hour)
		weather.callHook()
		return
	}
	//branch if status from client isn't OK or service unavailable and delete webhook
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		//send output to console
		fmt.Printf(
			"%v {\n" +
			"    id:              %s \n" +
			"    status_code:     %v,\n" +
			"    raw_error:       %s,\n" +
			"    message:         %s,\n" +
			"}\n", 
			time.Now().Format("2006-01-02 15:04:05"), weather.ID, res.StatusCode, "None",
			"Webhook URL is not valid. Deleting webhook...", 
		)
		err = storage.Firebase.Delete(dict.WEATHER_COLLECTION, weather.ID)
		if err != nil {
			//send output to console
			fmt.Printf(
				"%v {\n" +
				"    id:              %s \n" +
				"    status_code:     %v,\n" +
				"    raw_error:       %s,\n" +
				"    message:         %s,\n" +
				"}\n", 
				time.Now().Format("2006-01-02 15:04:05"), weather.ID, "None", err.Error(),
				"Didn't manage to delete webhook.", 
			)
		}
		dict.MutexState.Unlock()
		return
	}
	dict.MutexState.Unlock()
	weather.callHook()
}
