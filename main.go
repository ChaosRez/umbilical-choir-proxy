package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"math/big"
	"os"
	"time"
)

var (
	client  = resty.New()       // make a resty client
	host    = os.Getenv("HOST") // provided by agent in tf
	port    = os.Getenv("PORT")
	f1Name  = os.Getenv("F1NAME")
	f2Name  = os.Getenv("F2NAME")
	program = os.Getenv("PROGRAM")
	client     = resty.New().SetTimeout(500 * time.Millisecond) // make a resty client
)

func main() {
	startProxy := time.Now()
	// input from bash, at least contains an empty string + Env
	input := os.Args[1]

	// Call one of the function versions
	resp, err, elap, isF2 := uniformCallAndLog(input)
	if err != nil {
		log.Fatal("error running uniformCallAndLog:", err)
		return
	}

	// stdout
	fmt.Printf("resp: %s \n took: %s\n", resp, elap)

	// Push total proxy time
	elapTotal := time.Since(startProxy)

	payload := MetricUpdatePayload{
		Job:     "umbilical-choir",
		Program: program,
		Metrics: []Metric{
			{MetricName: "call_count", Value: 1},
			{MetricName: "proxy_time", Value: float64(elapTotal) / float64(time.Millisecond)},
		},
	}

	// append metrics
	if isF2 { // f2 was called and not f1
		newMetric1 := Metric{MetricName: "f2_count", Value: 1}
		newMetric2 := Metric{MetricName: "f2_time", Value: float64(elap) / float64(time.Millisecond)}
		payload.Metrics = append(payload.Metrics, newMetric1)
		payload.Metrics = append(payload.Metrics, newMetric2)
	} else {
		newMetric1 := Metric{MetricName: "f1_count", Value: 1}
		newMetric2 := Metric{MetricName: "f1_time", Value: float64(elap) / float64(time.Millisecond)}
		payload.Metrics = append(payload.Metrics, newMetric1)
		payload.Metrics = append(payload.Metrics, newMetric2)
	}

	// Push metrics to the aggregator
	errMetric := SendMetrics(host, 9999, payload)
	if errMetric != nil {
		log.Fatalf("Error sending metrics: %v\n", errMetric)
	}
}

// randomly calls one of the two functions. Returns the response, error, elapsed time, and a boolean indicating which function was called
func uniformCallAndLog(input string) (string, error, time.Duration, bool) {
	// Generate a random number (0 or 1)
	choice, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return "", err, 0, false // FIXME: this will be counted as f1 errors
	}

	if choice.Int64() == 0 { // call f1
		resp1, err1, elap1 := f1Call(input)
		if err1 != nil {
			fmt.Printf("error calling %s: %v\n", f1Name, err1)
		}

		return resp1, err1, elap1, false
	} else { // call f2
		resp2, err2, elap2 := f2Call(input)
		if err2 != nil {
			fmt.Printf("error calling %s: %v\n", f2Name, err2)
		}
		return resp2, err2, elap2, true
	}
}

func f1Call(input string) (string, error, time.Duration) {

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
	return checkResponse(call1Response)
}

func f2Call(input string) (string, error, time.Duration) {

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
	return checkResponse(call2Response)
}
