package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

func checkResponseAndReturnBody(fn func() (*resty.Response, error, time.Duration)) (string, error, time.Duration) {
	resp, err, runTime := fn()
	if err != nil {
		return "", err, runTime
	}
	if !resp.IsSuccess() { // if not in 2xx range
		msg := fmt.Sprintf("non-successful response (%d): %v", resp.StatusCode(), resp.String())
		return "", errors.New(msg), runTime
	}
	return string(resp.Body()), nil, runTime
}
