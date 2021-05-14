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
	arrWeather, err := storage.Firebase.GetAll(hookdb)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when initializing weather webhooks.\n\tRaw error: %v\n}\n",
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
		"%v {\n\tSuccesfully initialized weather webhooks.\n\tAmount: %v\n}\n",
		time.Now().Format("2006-01-02 15:04:05"), strconv.Itoa(len(arrWeather)),
	)
	diag.HookAmount += len(arrWeather)
}

// callHook calls webhook.
func (weather *Weather) callHook() {
	//check if webhook still exist in database
	_, exist := storage.Firebase.Get(hookdb, weather.ID)
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
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
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
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), recorder.Result().StatusCode,
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//convert from structure to bytes and branch if an error occurred
	output, err := json.Marshal(weatherDetails)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when parsing Weather structure.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//create new POST request and branch if an error occurred
	req, err = http.NewRequest(http.MethodPost, weather.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating new POST request.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//hash structure and branch if an error occurred
	mac := hmac.New(sha256.New, dict.Secret)
	_, err = mac.Write([]byte(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when hashing content before POST request.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
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
		fmt.Printf(
			"%v {\n\tError when sending HTTP content to webhook.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		dict.MutexState.Unlock()
		weather.callHook()
		return
	}
	//branch if status from client isn't OK or service unavailable and delete webhook
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		fmt.Printf(
			"%v {\n\tWebhook URL is not valid. Deleting webhook...\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), res.StatusCode,
		)
		err = storage.Firebase.Delete(hookdb, weather.ID)
		if err != nil {
			fmt.Printf(
				"%v {\n\tDidn't manage to delete webhook.\n\tRaw error: %v\n}\n",
				time.Now().Format("2006-01-02 15:04:05"), err.Error(),
			)
		}
		dict.MutexState.Unlock()
		return
	}
	dict.MutexState.Unlock()
	weather.callHook()
}
