package weatherHook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"main/api/weather"
	"main/db"
	"main/dict"
	"main/fun"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"

	"google.golang.org/api/iterator"
)

// We use mutex locks in callLoop to ensure we get no concurrency issues WRT map writes
// Since we're loading every hook in when the program starts running, any 2 or more hooks with the same timeout
// would otherwise be at risk of panicking when running at the same time.
var mutex = &sync.Mutex{}

/**
* InitHooks
* Initiates webhook triggers for all weather webhooks.
**/
func InitHooks(database *db.Database) error {
	iter := database.Client.Collection("weatherHookDB").Documents(database.Ctx)
	for {
        doc, err := iter.Next()
        if err == iterator.Done {
            break
		}
        if err != nil {
            return err
		}
		// Create dummy variables
		var temp WeatherHook
		var dbMap = doc.Data()
		// Extract data from iterator
		err = temp.ReadData(dbMap["Container"].(interface{}))
		temp.ID = doc.Ref.ID
		if err != nil {
			return err
		}
		// Start as go routine, else system will hang for the sleep time!
		go temp.callLoop()
	}
	//print message with amount of webhooks initilizied
	fmt.Printf(
		"%v {\n\tSuccesfully initialized WeatherHookwebhooks.\n}",
		time.Now().Format("2006-01-02 15:04:05"),
	)
	return nil
}


/**
* callLoop
* Function handling webhook triggering. It runs as a go routine every x hours, where x is the user-input timeout, for each webhook.
**/
func (weatherHook *WeatherHook) callLoop() {
	_, exist := db.DB.Get("weatherHookDB", weatherHook.ID)
	if !exist {
		return
	}	
	// Sleep
	fun.HookSleep(weatherHook.Timeout)
	// Lock and do its thing when the timeout is up
	mutex.Lock()

	url := dict.GetWeatherURL(weatherHook.Location, "")
	var weather weather.Weather
	//create new GET request and branch if an error occurred
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		mutex.Unlock()		
		weatherHook.callLoop()
		return
	}
	//call the policy handler and branch if the status code is not OK
	//this stops timed out request being sent to the webhook
	recorder := httptest.NewRecorder()
	weather.Handler(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		fmt.Printf(
			"%v {\n\tError when creating HTTP request to Weather.Handler().\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), recorder.Result().StatusCode,
		)
		mutex.Unlock()
		weatherHook.callLoop()
		return
	}
	//convert from structure to bytes and branch if an error occurred
	output, err := json.Marshal(weather)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when parsing Weather structure.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		mutex.Unlock()
		weatherHook.callLoop()
		return
	}
	//create new POST request and branch if an error occurred
	req, err = http.NewRequest(http.MethodPost, weatherHook.URL, bytes.NewBuffer(output))
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when creating new POST request.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		mutex.Unlock()
		weatherHook.callLoop()
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
		mutex.Unlock()
		weatherHook.callLoop()
		return
	}
	//convert hashed structure to string and add to header
	req.Header.Add("Signature", hex.EncodeToString(mac.Sum(nil)))
	//update header to JSON
	req.Header.Set("Content-Type", "application/json")
	//send request to client and branch if an error occured
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf(
			"%v {\n\tError when sending HTTP content to webhook.\n\tRaw error: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), err.Error(),
		)
		mutex.Unlock()
		weatherHook.callLoop()
		return
	}
	// Check URL is valid, delete webhook if not
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusServiceUnavailable {
		fmt.Printf(
			"%v {\n\tWebhook URL is not valid. Deleting webhook...\n\tStatus code: %v\n}\n",
			time.Now().Format("2006-01-02 15:04:05"), res.StatusCode,
		)
		err = db.DB.Delete("weatherHookDB", weatherHook.ID)
		mutex.Unlock()
		return
	} else {
	// When we've finished processing the trigger, we can unlock and recur in a new routine
		mutex.Unlock()
		weatherHook.callLoop()
	}
	return
}
