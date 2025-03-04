package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/big"
	"os"
	"strconv"
	"time"
)

var (
	client     = resty.New().SetTimeout(500 * time.Millisecond) // make a resty client
	f1Uri      = os.Getenv("F1ENDPOINT")                        // provided by agent in tf
	f2Uri      = os.Getenv("F2ENDPOINT")
	agentHost  = os.Getenv("AGENTHOST")
	f1Name     = os.Getenv("F1NAME")
	f2Name     = os.Getenv("F2NAME")
	program    = os.Getenv("PROGRAM")
	bChanceStr = os.Getenv("BCHANCE")
	bChance    int
)

func fn(data string, headers map[string]string) (string, error) {
	startProxy := time.Now()
	var err error
	bChance, err = strconv.Atoi(bChanceStr)
	if err != nil {
		return "", fmt.Errorf("error converting BCHANCE '%v' to int: %v", bChanceStr, err)
	}

	fmt.Printf("Call with input: %s\n", data)

	input := data
	//functionChoice := ""
	functionChoice := headers["X-Function-Choice"]

	var resp string
	var funcCallLatency time.Duration
	var isF2 bool

	// Check the function choice
	switch functionChoice {
	case "f1":
		resp, err, funcCallLatency = f1Call(input)
		isF2 = false
	case "f2":
		resp, err, funcCallLatency = f2Call(input)
		isF2 = true
	default:
		// Randomly call one of the function versions based on bChance
		resp, err, funcCallLatency, isF2 = randomCallAndLog(input, bChance)
	}

	// stdout
	//if err == nil {
	//	fmt.Printf("resp: %s \n took: %s\n", resp, funcCallLatency)
	//}

	// Push total proxy time
	elapTotal := time.Since(startProxy)

	payload := MetricUpdatePayload{
		Program: program,
		Metrics: []Metric{
			{MetricName: "call_count", Value: 1},
			{MetricName: "proxy_time", Value: float64(elapTotal) / float64(time.Millisecond)},
		},
	}

	// append metrics
	if isF2 { // f2 was called and not f1
		newMetric1 := Metric{MetricName: "f2_count", Value: 1}
		payload.Metrics = append(payload.Metrics, newMetric1)
		if err != nil { // error running randomCallAndLog
			newMetric2 := Metric{MetricName: "f2_error_count", Value: 1}
			payload.Metrics = append(payload.Metrics, newMetric2)
			fmt.Errorf("error running f2:", err)
		} else {
			newMetric2 := Metric{MetricName: "f2_time", Value: float64(funcCallLatency) / float64(time.Millisecond)}
			payload.Metrics = append(payload.Metrics, newMetric2)
		}
	} else {
		newMetric1 := Metric{MetricName: "f1_count", Value: 1}
		payload.Metrics = append(payload.Metrics, newMetric1)
		if err != nil { // error running randomCallAndLog
			newMetric2 := Metric{MetricName: "f1_error_count", Value: 1}
			payload.Metrics = append(payload.Metrics, newMetric2)
			return "", fmt.Errorf("error running f1:", err)
		} else {
			newMetric2 := Metric{MetricName: "f1_time", Value: float64(funcCallLatency) / float64(time.Millisecond)}
			payload.Metrics = append(payload.Metrics, newMetric2)
		}
	}

	// Push metrics to the aggregator
	startMetric := time.Now()
	errMetric := SendMetrics(agentHost, 9999, payload)
	elapMetric := time.Since(startMetric)

	if errMetric != nil {
		return "", fmt.Errorf("error sending metrics: %v", errMetric)
	}

	out := fmt.Sprintf("resp: %s \nfunc call took: %s\nSendMetrics took: %s", resp, funcCallLatency, elapMetric)
	return out, nil
}

// randomly calls one of the two functions. Returns the response, error, elapsed time, and a boolean indicating which function was called
func randomCallAndLog(input string, bChance int) (string, error, time.Duration, bool) {
	// Generate a random number between 0 and 100
	choice, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return "", err, 0, false // FIXME: this will be counted as f1 errors
	}

	// Choose a function version to call
	if choice.Int64() < int64(bChance) { // call f2
		resp2, err2, elap2 := f2Call(input)
		if err2 != nil {
			fmt.Printf("error calling %s: %v\n", f2Name, err2)
		}
		return resp2, err2, elap2, true

	} else { // call f1
		resp1, err1, elap1 := f1Call(input)
		if err1 != nil {
			fmt.Printf("error calling %s: %v\n", f1Name, err1)
		}
		return resp1, err1, elap1, false
	}
}

func f1Call(input string) (string, error, time.Duration) {

	call1Response := func() (*resty.Response, error, time.Duration) {
		start1 := time.Now()
		resp1, err1 := client.R().
			//EnableTrace().
			SetBody(input).
			Post(f1Uri)
		elap1 := time.Since(start1)
		if err1 != nil {
			return nil, err1, elap1
		}
		return resp1, nil, elap1
	}

	// validate the response
	return checkResponseAndReturnBody(call1Response)
}

func f2Call(input string) (string, error, time.Duration) {

	call2Response := func() (*resty.Response, error, time.Duration) {
		start2 := time.Now()
		resp2, err2 := client.R().
			//EnableTrace().
			SetBody(input).
			Post(f2Uri)
		elap2 := time.Since(start2)
		if err2 != nil {
			return nil, err2, elap2
		}
		return resp2, nil, elap2
	}

	// validate the response
	return checkResponseAndReturnBody(call2Response)
}

// wiki:
//// Access the trace information
//traceInfo := resp.Request.TraceInfo()
//fmt.Printf("DNS Lookup: %v\n", traceInfo.DNSLookup)
//fmt.Printf("TCP Connection: %v\n", traceInfo.ConnTime)
//fmt.Printf("TLS Handshake: %v\n", traceInfo.TLSHandshake)
//fmt.Printf("Server Processing: %v\n", traceInfo.ServerTime)
//fmt.Printf("Response Time: %v\n", traceInfo.ResponseTime)
//fmt.Printf("Total Time: %v\n", traceInfo.TotalTime)
