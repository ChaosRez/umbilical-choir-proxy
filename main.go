package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"time"
)

func main() {
	// input from bash, at least contains an empty string + Env
	input := os.Args[1]
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	f1Name := os.Getenv("F1NAME")
	f2Name := os.Getenv("F2NAME")

	// make a resty client
	client := resty.New()

	// call f1
	call1Response := func() (*resty.Response, error, time.Duration) {
		start1 := time.Now()
		resp1, err1 := client.R().
			EnableTrace().
			SetBody(input).
			Post(fmt.Sprintf("http://%s:%s/%s", host, port, f1Name))
		elap1 := time.Since(start1)
		if err1 != nil {
			return nil, err1, elap1
		}
		return resp1, nil, elap1
	}

	// validate the response
	resp1, err1, elap1 := checkResponse(call1Response)
	if err1 != nil {
		fmt.Printf("error calling %s: %v\n", f1Name, err1)
		return
	}

	// call f2
	call2Response := func() (*resty.Response, error, time.Duration) {
		start2 := time.Now()
		resp2, err2 := client.R().
			EnableTrace().
			SetBody(input).
			Post(fmt.Sprintf("http://%s:%s/%s", host, port, f2Name))
		elap2 := time.Since(start2)
		if err2 != nil {
			return nil, err2, elap2
		}
		return resp2, nil, elap2
	}

	// validate the response
	resp2, err, elap2 := checkResponse(call2Response)
	if err != nil {
		fmt.Printf("error calling %s: %v\n", f2Name, err)
		return
	}

	// stdout
	fmt.Printf("%s took %v and %s took %v\n", f1Name, elap1, f2Name, elap2)
	fmt.Printf("resp1: %s \n resp2: %s\n", resp1, resp2)
}

func checkResponse(fn func() (*resty.Response, error, time.Duration)) (string, error, time.Duration) {
	resp, err, runTime := fn()
	if err != nil {
		return "", err, runTime
	}
	if !resp.IsSuccess() {
		msg := fmt.Sprintf("non-successful response (%d)", resp.StatusCode())
		return "", errors.New(msg), runTime
	}
	return string(resp.Body()), nil, runTime
}
