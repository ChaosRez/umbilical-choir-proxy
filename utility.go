package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"time"
)

// Create a new Pusher targeting the Pushgateway
var pusher = push.New(fmt.Sprintf("http://%s:%s", host, "9091"), "umbilical-choir")

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

//func PushResponseTime(metric prometheus.Gauge, duration time.Duration) error {
//	// Update the metric locally
//	metric.Set(float64(duration) / float64(time.Millisecond))
//
//	// Set job and groupings (optional)
//	pusher.Grouping("program", program)
//
//	return pusher.Collector(metric).Add()
//}

func PushMetrics(metrics ...prometheus.Gauge) error {
	// Set job and groupings (optional)
	pusher.Grouping("program", program)

	for _, metric := range metrics {
		pusher.Collector(metric)
	}
	return pusher.Add()
}
