package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// requestData gets raw data from REST services based on URL.
func RequestData(url string) ([]byte, int, error) {
	//create new request and branch if an error occurred
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, req.Response.StatusCode, err
	}
	//timeout after 2 seconds
	apiClient := http.Client{
		Timeout: time.Second * 4,
	}
	//get response and branch if an error occurred
	res, err := apiClient.Do(req)
	if err != nil {
		return nil, http.StatusRequestTimeout, err
	}
	if res.StatusCode != http.StatusOK {
		err = errors.New("requesting data: status code is not OK")
		return nil, res.StatusCode, err
	}
	//read output and branch if an error occurred
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, http.StatusRequestTimeout, err
	}
	return output, http.StatusOK, nil
}
