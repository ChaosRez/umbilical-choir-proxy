package main

import (
	"crypto/rand"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"math/big"
	"os"
	"time"
)

var (
	responseTime1 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "response_time_1",
		Help: "Response time for the f1 (ms)",
	})
	responseTime2 = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "response_time_2",
		Help: "Response time for the f2 (ms)",
	})
	proxyTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "proxy_time",
		Help: "Total proxy runtime (ms)",
	})
	client  = resty.New()       // make a resty client
	host    = os.Getenv("HOST") // provided by agent in tf
	port    = os.Getenv("PORT")
	f1Name  = os.Getenv("F1NAME")
	f2Name  = os.Getenv("F2NAME")
	program = os.Getenv("PROGRAM")
)

func main() {
	startProxy := time.Now()
	// input from bash, at least contains an empty string + Env
	input := os.Args[1]

	// Call one of the function versions
	resp, err, elap := uniformCallAndLog(input)

	// stdout
	fmt.Printf("resp: %s \n took: %s\n", resp, elap)

	// Push total proxy time
	elapTotal := time.Since(startProxy)
	if err = PushResponseTime(proxyTime, elapTotal); err != nil {
		fmt.Println("Failed to push proxy_time to Pushgateway:", err)
	}
}

// randomly calls one of the two functions
func uniformCallAndLog(input string) (string, error, time.Duration) {
	// Generate a random number (0 or 1)
	choice, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		fmt.Println("Error generating random number:", err)
		return "", err, 0
	}

	if choice.Int64() == 0 { // call f1
		resp1, err1, elap1 := f1Call(input)
		if err1 != nil {
			fmt.Printf("error calling %s: %v\n", f1Name, err1)
		}

		// Push the updated metric to Pushgateway
		if err = PushResponseTime(responseTime1, elap1); err != nil {
			fmt.Println("Failed to push response_time_1 to Pushgateway:", err)
		}
		return resp1, err1, elap1
	} else { // call f2
		resp2, err2, elap2 := f2Call(input)
		if err2 != nil {
			fmt.Printf("error calling %s: %v\n", f2Name, err2)
		}

		// Push the updated metric to Pushgateway
		if err = PushResponseTime(responseTime2, elap2); err != nil {
			fmt.Println("Failed to push response_time_2 to Pushgateway:", err)
		}

		return resp2, err2, elap2
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

func init() {
	prometheus.MustRegister(responseTime1)
	prometheus.MustRegister(responseTime2)
}
