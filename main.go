package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
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
	host    = os.Getenv("HOST")
	port    = os.Getenv("PORT")
	f1Name  = os.Getenv("F1NAME")
	f2Name  = os.Getenv("F2NAME")
	program = os.Getenv("PROGRAM")
)

func main() {
	// input from bash, at least contains an empty string + Env
	input := os.Args[1]

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

	// Push the updated metric to Pushgateway
	if err = PushResponseTime(responseTime1, elap1); err != nil {
		fmt.Println("Failed to push response_time_1 to Pushgateway:", err)
	}

	if err = PushResponseTime(responseTime2, elap2); err != nil {
		fmt.Println("Failed to push response_time_2 to Pushgateway:", err)
		//return
	}
}

func init() {
	prometheus.MustRegister(responseTime1)
	prometheus.MustRegister(responseTime2)
}
