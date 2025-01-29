package http

import (
	"net/http"
	"time"
)

const maxRetries = 3
const initialDelay = 2 * time.Second

func DoRequest(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {
		resp, err = http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}
		time.Sleep(initialDelay * time.Duration(i+1))
	}
	return resp, err
}
