package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Get will send a get request
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status Code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
